import React, { memo, useEffect, useState } from "react";
import { Spin } from "@douyinfe/semi-ui";
import { IconLink } from "@douyinfe/semi-icons";
import type { BriefWorkItem } from "@lark-project/js-sdk";
import { fetchSmartElfConfig } from "../../api/services";
import sdk from "../../utils/sdk";
import "./index.less";

export const DisplayComp = memo(() => {
  const [projectKey, setProjectKey] = useState<string>("");
  const [activeWorkItem, setActiveWorkItem] = useState<
    BriefWorkItem | undefined
  >(undefined);
  const [filedId, setFiledId] = useState("");
  const [userInfo, setUserInfo] = useState<{ name: string; openId: string }>();
  const [spinning, setSpinning] = useState(false);
  useEffect(() => {
    sdk.Context.load().then(({ mainSpace, activeWorkItem }) => {
      mainSpace?.id && setProjectKey(mainSpace?.id);
      setActiveWorkItem(activeWorkItem);
    });
  }, []);
  useEffect(() => {
    setSpinning(true);
    if (projectKey && activeWorkItem && filedId) {
      const { workObjectId, id: workItemId } = activeWorkItem;
      sdk.WorkItem.load({
        spaceId: projectKey,
        workObjectId,
        workItemId,
      })
        .then(async (workItem) => {
          const val = (await workItem.getFieldValue(filedId))?.split("###");
          setUserInfo({
            name: val?.[0],
            openId: val?.[1],
          });
        })
        .catch(() => {
          setUserInfo({
            name: "",
            openId: "",
          });
        })
        .finally(() => {
          setSpinning(false);
        });
    } else {
      setSpinning(false);
    }
  }, [activeWorkItem, filedId, projectKey]);
  useEffect(() => {
    if (!projectKey) {
      return;
    }
    fetchSmartElfConfig(projectKey).then((res) => {
      setFiledId(res?.config?.creator_field_key);
    });
  }, [projectKey]);
  return (
    <Spin tip="加载中..." spinning={spinning}>
      <div className="user-select--tag">
        {userInfo?.name ? (
          <>
            <span className="user-select--tag-name">{userInfo?.name}</span>
            <a
              target="_blank"
              rel="noreferrer"
              onClick={() => {
                sdk.navigation.open(`https://applink.feishu.cn/client/chat/open?openId=${userInfo?.openId}`, '_blank');
              }}
            >
            <IconLink className="user-select--tag-icon" />
            </a>
          </>
        ) : (
          "暂无数据"
        )}
      </div>
    </Spin>
  );
});
