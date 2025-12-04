import { snakeCase } from 'lodash';
import sdk from './sdk';

export const handleErrorMsg = (e: any, minVersion?: string) => {
  let msg = '';
  if (e.name === 'NotSupportedError') {
    msg = minVersion
      ? `当前客户端暂不支持，\n请升级飞书客户端到${minVersion}及以上版本`
      : '当前客户端暂不支持，\n请升级飞书客户端到最新版本';
  } else {
    msg = '内部错误:' + (e.message || e.originMessage);
  }
  document.body.appendChild(document.createTextNode(msg));
};

export const getLang = async () => {
  const { language } =
    (await sdk.Context.load().catch((e) => handleErrorMsg(e))) || {};
  return language || 'zh_CN';
};

export const getStorage = (key: string) =>
  sdk.storage
    .getItem(key)
    .then((res) => res ?? null)
    .catch((e) => handleErrorMsg(e));

export const setStorage = (key: string, value?: string) => {
  sdk.storage.setItem(key, value).catch((e) => handleErrorMsg(e));
};

export const removeStorage = (key: string) =>
  sdk.storage.removeItem(key).catch((e) => handleErrorMsg(e));

export const getProjectKey = async () => {
  const context = await sdk.Context.load().catch((e) => handleErrorMsg(e));
  return context?.mainSpace?.id || '';
};

export const getUserKey = async (): Promise<string> => {
  const context = await sdk.Context.load().catch((e) => handleErrorMsg(e));
  return context?.loginUser.id || '';
};

export const getSpace = (projectKey: string) =>
  sdk.Space.load(projectKey).catch((e) => handleErrorMsg(e));

export const getAllWorkObjectList = async (projectKey: string) => {
  const space = await getSpace(projectKey);
  if (space) {
    return space.allWorkObjectList;
  }
  return [];
};

export const getControlContext = async () =>
  sdk.control.getContext().catch((e) => handleErrorMsg(e, '7.25.0'));

/**
 * 判断控件在详情页还是节点表单
 * @param tab
 * 1.内置的tab是内置的名字 detail、comment自定义tab是随机生成的
 * 2.节点的 tab 是拼的。空间:工作项:模版 uuid:节点 key: form_conf
 * @returns boolen
 */
export const checkIsNodeForm = (tab: string) =>
  /:.*:(\b(form_conf)\b)/g.test(tab);

interface FieldOption {
  label: string;
  value: string | number;
  children?: FieldOption[];
  display?: string;
}

export const getMapByLists = (lists: FieldOption[] | undefined) => {
  if (!lists) {
    return;
  }
  const listMaps: Record<string, any> = {};
  // 递归遍历业务线并设置display和detail
  const recursiveMap = (item: FieldOption, itemParent?: FieldOption) => {
    const display = itemParent
      ? `${itemParent.display || itemParent.label}/${item.label}`
      : item.label;
    item.display = display;
    listMaps[item.value] = {
      value: item.value,
      label: item.label,
      display,
    };
    // children是可选字段
    if ('children' in item) {
      item.children?.forEach((child) => {
        recursiveMap(child, item);
      });
    }
  };
  lists.forEach((item) => {
    recursiveMap(item);
  });
  return listMaps;
};

export const getLarkDefinitionMap = (value, key = 'id') =>
  value.reduce((accumulator, current) => {
    accumulator[current[key]] = current;
    return accumulator;
  }, {});

function isObject(value: any) {
  return Object.prototype.toString.call(value) === '[object Object]';
}

const convertToSnake = (str: string) => snakeCase(str);

export const keysToSnakeCase = (obj: Record<string, any>): any => {
  if (!isObject(obj)) {
    return obj;
  }
  if (Array.isArray(obj)) {
    return obj.map((item) => keysToSnakeCase(item));
  }
  const result = {};
  for (const key of Object.keys(obj)) {
    let val = obj[key];
    if (Array.isArray(val)) {
      val = val.map((item) => keysToSnakeCase(item));
    } else if (isObject(val)) {
      val = keysToSnakeCase(val);
    }
    const snakeKey = convertToSnake(key);
    result[snakeKey] = val;
  }
  return result;
};

export function recursiveJSONParse(input) {
  // 尝试解析输入
  try {
    const result = JSON.parse(input);
    // 如果结果仍然是字符串，尝试递归解析
    if (typeof result === 'string') {
      return recursiveJSONParse(result);
    } else {
      // 否则返回结果
      return result;
    }
  } catch (error) {
    // 如果解析失败，抛出错误
    throw new Error(`Failed to parse JSON: ${error}`);
  }
}
