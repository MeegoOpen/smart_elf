import axios from 'axios';
import { apiHost } from '../constants/index';
import { getLang, getUserKey } from '../utils/index';


// 创建 axios 实例
const request = axios;

// 请求拦截器
request.interceptors.request.use(
  async (config) => {
    // 在请求发送之前做一些处理
    // 添加请求头信息
    if (config.url?.startsWith('/')) {
      config.url = apiHost + config.url;
    }
    const lang = await getLang();
    config.headers['X-USER-KEY'] = await getUserKey();
    config.headers.locale = lang;

    return config;
  },
  (error) => {
    // 对请求错误做些什么
    console.error(error);
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  async (response) => {
    // 在响应之前做一些处理
    const res = response.data;
    if (res.err_code !== 0) {
      console.error(
        `[Request error] url=${response.config.url}, code=${res.err_code}, message=${res.err_msg}, x-tt-logid=${response.headers['x-tt-logid']}`
      );
    }
    return res;
  },
  (error) => {
    // 对响应错误做些什么
    // toastCallBack(error);
    return Promise.reject(error);
  }
);

class Semaphore {
  private limit: number;
  private count: number; // 计数
  private lastRunTime: number; // 计时
  private duration = 1000;

  constructor(limit: number) {
    this.limit = limit;
    this.count = 0;
    this.lastRunTime = Date.now();
  }

  public acquire() {
    // 每次发起都判断下是否超过了 1s，如果超过就更新下计数和计时
    if (Date.now() - this.lastRunTime > this.duration) {
      this.lastRunTime = Date.now();
      this.count = 0;
    }
    // 靠 Promise 控制 resolve 的触发来实现阻塞
    return new Promise<void>((resolve) => {
      // 计数未超过限制时正常执行
      if (this.count < this.limit) {
        this.count++;
        resolve();
      } else {
        // 计数超过限制后开始无限循环
        const timer = setInterval(() => {
          if (Date.now() - this.lastRunTime > this.duration) {
            // 直到超过 1s 后，更新计数和计时
            this.lastRunTime = Date.now();
            this.count = 0;
          }
          if (this.count < this.limit) {
            // 满足条件时，正常执行，跳出循环
            this.count++;
            resolve();
            clearInterval(timer);
          }
        }, 100);
      }
    });
  }
}

type AsyncFunction = (...args: any) => Promise<any>;

const limitQps = (fn: AsyncFunction, limit: number): AsyncFunction => {
  const semaphore = new Semaphore(limit);

  return async (...args): Promise<any> => {
    await semaphore.acquire();

    try {
      return await fn(...args);
    } catch (e) {
      console.error('limitQps', e);
    }
  };
};

// 定义 requestCache 装饰器工厂函数
export const requestCache = (
  originalFunction: (...args: any[]) => Promise<any>,
  cacheDuration: number
) => {
  const cache: Map<string, any> = new Map();
  const cachePromise: Map<string, Promise<any>> = new Map();

  // 返回一个装饰器函数
  return (...args: any[]) => {
    const [url, payload] = args;

    // 生成唯一的缓存键
    const cacheKey = `${url}-${JSON.stringify(payload)}`;

    // 检查是否有缓存
    if (cache.has(cacheKey)) {
      return Promise.resolve(cache.get(cacheKey));
    }

    // 检查是否有进行中的 Promise
    if (!cachePromise.has(cacheKey)) {
      cachePromise.set(
        cacheKey,
        originalFunction(...args)
          .then((result) => {
            // 请求完成后缓存结果
            cache.set(cacheKey, result);
            // 缓存过期后清除缓存，防止内存泄露
            setTimeout(() => {
              cache.delete(cacheKey);
            }, cacheDuration);

            return result;
          })
          .finally(() => {
            // 请求结束后清除 Promise 缓存
            cachePromise.delete(cacheKey);
          })
      );
    }
    return cachePromise.get(cacheKey)!;
  };
};

export const limitPost: typeof request.post = requestCache(
  limitQps(request.post, 9),
  5000
);

export const limitGet: typeof request.get = requestCache(
  limitQps(request.get, 9),
  5000
);

export default request;

export const mockFetch = <D = any>(data: D, ms = 500) =>
  new Promise((rev) => {
    setTimeout(() => rev(data), ms);
  });
