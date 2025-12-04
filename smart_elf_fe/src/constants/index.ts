import manifestConfig from '../../plugin.config.json';

export const { pluginId: APP_KEY, siteDomain, OpenApiHost } = manifestConfig;

export const apiHost = OpenApiHost;
export const requestHost = siteDomain;

export const filedBlackList = ['abort_detail', 'name'];
