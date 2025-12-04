import { filedBlackList } from "../constants";
import request from "./request";
import { apiHost } from "../constants";
interface ResponseWrap<D> {
  err_code: number;
  err_msg: string;
  data?: D;
  error?: {
    id: number;
    localizedMessage: {
      locale: string;
      message: string;
    };
  };
}



export interface ISmartElf {
  project_key: string;
  config: {
    bot_info: {
      bot_id: string;
      bot_secret: string;
      verification_token: string;
    };
    reply_switch: boolean;
    work_item_type_key: string;
    creator_field_key: string;
    create_group_switch: boolean;
    api_user_key: string;
    work_item_template_id: number;
    work_item_api_name: string;
  };
}

export const fetchSmartElfConfig = (project_key: string) =>
  request
    .get<Omit<ISmartElf, "project_key">>(
    `${apiHost}/api/v1/config/query?project_key=${project_key}`
    )
    .then((res) => res.data);

export const updateSmartElfConfig = (param: ISmartElf) =>
  request.post<
    unknown,
    ResponseWrap<{
      app_type: string;
      expire_time: number;
      token: string;
    }>
    >(`${apiHost}/api/v1/config/update`, {
    ...param,
  });

export interface WorkObjectField {
  field_alias?: string;
  key: string;
  name: string;
  type: string;
  field_key: string;
  field_name: string;
  field_type_key: string;
  is_custom_field?: boolean;
  is_obsoleted?: boolean;
  value_generate_mode: number | string;
}

export const fetchWorkObjectFields = (
  projectKey: string,
  workItemKey: string
) =>
  request
    .post<unknown, ResponseWrap<WorkObjectField[]>>(
     `${apiHost}/proxy/open_api/${projectKey}/field/all`,
      {
        work_item_type_key: workItemKey,
      }
    )
    .then(({ err_code, data }) => {
      let curData: WorkObjectField[] = [];
      if (Array.isArray(data)) {
        curData = data
          .filter((fd) => !filedBlackList.includes(fd.field_key))
          ;
      }
      return {
        err_code: err_code,
        data: curData,
      };
    });

export const fetchSmartElfSig = (project_key: string) =>
  request
    .post<{ signature: string }>(
    `${apiHost}/api/v1/config/signature`,{
        project_key,
    }
    )
    .then((res) => res.data);
