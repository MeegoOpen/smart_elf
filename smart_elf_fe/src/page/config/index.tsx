import React, { useEffect, useRef, useState } from "react";
import {
  Form,
  Button,
  Toast,
  Card,
  Skeleton,
} from "@douyinfe/semi-ui";
import _ from "lodash";
import sdk from "../../utils/sdk";
import {
  fetchSmartElfConfig,
  updateSmartElfConfig,
  fetchWorkObjectFields,
  fetchSmartElfSig,
  type ISmartElf,
} from "../../api/services";
import { apiHost } from "../../constants";
window.JSSDK.utils.overwriteThemeForSemiUI();

export const getSpace = (projectKey: string) => sdk.Space.load(projectKey);

const getRoleList = async (projectKey: string, workItemKey: string) => {
  const { err_code, data } = await fetchWorkObjectFields(
    projectKey,
    workItemKey
  );
  if (err_code === 0) {
    return _.uniqBy(
      data
        .filter(
          (i) =>
            i.field_type_key === "text" && i.value_generate_mode !== "Calculate"
        )
        .map(({ field_name, field_key }) => ({
          value: field_key,
          label: field_name,
        })),
      "value"
    );
  }
  return [];
};
const getWorkItemTempleteList = async (
  projectKey: string,
  workItemKey: string
) => {
  const detail = await sdk.WorkObject.load({
    spaceId: projectKey,
    workObjectId: workItemKey,
  });
  const templateList = await detail.getTemplateList();
  return templateList
    .filter(({ disabled }) => disabled === false)
    .map(({ id, name }) => ({
      value: id,
      label: name,
    }));
};
const INIT_VALUES: Partial<ISmartElf["config"]> = {
  reply_switch: false,
  create_group_switch: false,
};
const config = () => {
  const formApiRef = useRef<any>();
  const [projectKey, setProjectKey] = useState<string>("");
  const [workItemList, setWorkItemList] = useState<
    { value: string; label: string; apiName: string }[]
  >([]);
  const [workItemTempleteList, setWorkItemTempleteList] = useState<
    { value: number; label: string }[]
  >([]);
  const [roleList, setRoleList] = useState<{ value: string; label: string }[]>(
    []
  );
  const [checkIsBot, setCheckIsBot] = useState(false);
  const [disabled, setDisabled] = useState(true);
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    sdk.Context.load().then(({ mainSpace }) => {
      mainSpace?.id && setProjectKey(mainSpace?.id);
    });
  }, []);
  useEffect(() => {
    if (!projectKey) {
      return;
    }
    getSpace(projectKey).then(({ enabledWorkObjectList }) => {
      setWorkItemList(
        enabledWorkObjectList.map(({ id, name, apiName }) => ({
          value: id,
          label: name,
          apiName,
        }))
      );
    });
  }, [projectKey]);
  useEffect(() => {
    if (!projectKey || !formApiRef.current) {
      return;
    }
    fetchSmartElfConfig(projectKey)
      .then(async (res) => {
        const {
          bot_info: { bot_id, bot_secret, verification_token },
          reply_switch,
          work_item_type_key,
          creator_field_key,
          create_group_switch,
          api_user_key,
          work_item_template_id,
          work_item_api_name,
        } = res?.config;
        const formApi = formApiRef.current;
        formApi.setValues(
          {
            bot_id,
            bot_secret,
            verification_token,
            reply_switch,
            work_item_type_key,
            creator_field_key,
            create_group_switch,
            api_user_key,
            work_item_template_id:
              work_item_template_id === 0 ? undefined : work_item_template_id,
            work_item_api_name,
          },
          { isOverride: true }
        );
        getRoleList(projectKey, work_item_type_key).then(setRoleList);
        getWorkItemTempleteList(projectKey, work_item_type_key).then(
          setWorkItemTempleteList
        );
        setCheckIsBot(reply_switch || create_group_switch);
      })
      .catch((err) => {
        console.error("err", err);
      })
      .finally(() => setLoading(false));
  }, [projectKey, formApiRef.current]);
  const handleSubmit = () => {
    formApiRef.current
      .validate()
      .then((values) => {
        if (!projectKey) {
          return;
        }
        const {
          bot_id = "",
          bot_secret = "",
          verification_token,
          reply_switch,
          work_item_type_key,
          creator_field_key,
          create_group_switch,
          api_user_key,
          work_item_template_id,
          work_item_api_name,
        } = values;
        updateSmartElfConfig({
          project_key: projectKey,
          config: {
            bot_info: {
              bot_id,
              bot_secret,
              verification_token,
            },
            reply_switch,
            work_item_type_key,
            creator_field_key,
            create_group_switch,
            api_user_key,
            work_item_template_id,
            work_item_api_name,
          },
        }).then(({ err_code }) => {
          if (err_code === 0) {
            setDisabled(true);
            return Toast.info({ content: "已提交" });
          }
        });
      })
      .catch((errors) => {
        console.log(errors);
      });
  };
  const handleCopy = async () => {
    const { signature } = await fetchSmartElfSig(projectKey);

    if (signature) {
      const href = await sdk.navigation.getHref();
      const url = new URL(href);
      const success = await sdk.clipboard.writeText(
        `${apiHost}/api/v1/lark/event?sig=${signature}`
      );
      return success
        ? Toast.success({ content: "已复制webhook" })
        : Toast.error({ content: "复制失败，请重新复制webhook" });
    }
    return Toast.error({ content: "复制失败，请重新复制webhook" });
  };
  const workItemHandle = (val: string) => {
    const { setValue, getValue } = formApiRef.current;
    if (getValue("work_item_type_key") !== val) {
      setRoleList([]);
      setWorkItemTempleteList([]);
      setValue("creator_field_key", undefined);
      setValue("work_item_template_id", undefined);
      const apiName = workItemList.find((i) => i.value === val)?.apiName;
      setValue("work_item_api_name", apiName);
      getRoleList(projectKey, val).then(setRoleList);
      getWorkItemTempleteList(projectKey, val).then(setWorkItemTempleteList);
    }
  };
  return (
    <Form
      initValues={INIT_VALUES}
      getFormApi={(api) => (formApiRef.current = api)}
      style={{ paddingLeft: 24, margin: "0 auto" }}
      disabled={disabled}
      labelPosition="left"
      labelWidth="180px"
    >
        <Card
          shadows="always"
          style={{
            margin: "0 auto",
            maxWidth: "900px",
            padding: "16px 24px 32px",
          }}
          title="小精灵配置"
          headerExtraContent={
            <>
              <Button theme="borderless" onClick={handleCopy}>
                复制webhook
              </Button>
              {disabled ? (
                <Button
                  style={{ marginLeft: "10px" }}
                  onClick={() => setDisabled(false)}
                >
                  编辑
                </Button>
              ) : (
                <Button
                  style={{ marginLeft: "10px" }}
                  type="primary"
                  onClick={handleSubmit}
                >
                  保存
                </Button>
              )}
            </>
          }
        >
          
          <Skeleton
            placeholder={<Skeleton.Paragraph rows={2} />}
            loading={loading}
            active={true}
          >
            <Card title="工单配置" style={{ marginBottom: 20 }}>
              <Form.Select
                field="work_item_type_key"
                label="工作项"
                placeholder="请选择工作项"
                rules={[{ required: true, message: "必填" }]}
                style={{ width: "100%" }}
                optionList={workItemList}
                onChange={workItemHandle}
              />
              <Form.Select
                field="work_item_template_id"
                label="工作项模板"
                placeholder="请选择工作项模板"
                rules={[{ required: true, message: "必填" }]}
                style={{ width: "100%" }}
                optionList={workItemTempleteList}
              />
              <Form.Select
                field="creator_field_key"
                label="提交人存储字段"
                placeholder="请选择提交人存储字段"
                rules={[{ required: true, message: "必填" }]}
                style={{ width: "100%" }}
                optionList={roleList}
              />
              <Form.Input
                field="api_user_key"
                label="User Key"
                placeholder="请输入userkey"
                rules={[{ required: true, message: "必填" }]}
              />
            </Card>
            <Card title="飞书机器人配置" style={{ marginBottom: 20 }}>
              <Form.Input
                field="bot_id"
                label="App Id"
                style={{ width: "100%" }}
                rules={[{ required: checkIsBot, message: "必填" }]}
                placeholder="请输入 bot_id"
              ></Form.Input>
              <Form.Input
                field="bot_secret"
                label="App Secret"
                style={{ width: "100%" }}
                rules={[{ required: checkIsBot, message: "必填" }]}
                placeholder="请输入 bot_secret"
              ></Form.Input>

              <Form.Input
                field="verification_token"
                label="Verification Token "
                placeholder="请输入Verification Token "
                rules={[{ required: true, message: "必填" }]}
              />
            </Card>
            <Card title="操作配置" style={{ marginBottom: 20 }}>
              <Form.Switch
                label="是否开启自动回复"
                field="reply_switch"
                onChange={setCheckIsBot}
              />
              <Form.Switch
              label="是否自动创建群组"
              field="create_group_switch"
              onChange={setCheckIsBot}
            />
            </Card>
          </Skeleton>
        </Card>
    </Form>
  );
};

export default config;
