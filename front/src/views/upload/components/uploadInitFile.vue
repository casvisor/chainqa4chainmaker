<template>
  <t-card>
    <t-config-provider :global-config="languageConfig">
      <t-upload
        v-model="files"
        accept=".xls,.xlsx"
        placeholder="最多批量转换5个文件，重名文件不重复转换，支持xls、xlsx格式，文件每行的列数必须相等，单元格中的空值转换时会以_代替"
        theme="file-flow"
        multiple
        :disabled="onUploadLoading"
        :abridge-name="[25, 7]"
        :auto-upload="false"
        :max="5"
        tips="本部分处理仅在前端浏览器完成，将Excel文件转换为适宜IPFS存储的特定纯文本格式，转换至少一个文件后请点击下一步"
        :show-thumbnail="true"
        :allow-upload-duplicate-file="false"
        :is-batch-upload="false"
        :request-method="handleUpload"
        :cancelUploadButton="null"
        :showUploadProgress="false"
        :onRemove="removeFile"
      ></t-upload>
    </t-config-provider>

    <!--  -->
    <div style="margin-top: 20px" v-if="successFileDataList.length > 0">
      <t-typography-text strong>
        您已经成功转换了{{
          successFileDataList.length
        }}个文件，点击下方按钮进行确认并上载IPFS:
      </t-typography-text>
      <t-button
        block
        theme="primary"
        variant="outline"
        shape="round"
        style="margin-top: 10px"
        @click="handleNextStep"
      >
        下一步</t-button
      >
    </div>
  </t-card>
</template>

<script setup>
import { reactive, ref, defineEmits } from "vue";
import * as XLSX from "xlsx";

const files = ref([]);

const uploadMethod = ref("requestSuccessMethod");
const onUploadLoading = ref(false);

/**
 *
 * @param file_ 上传的文件
 * 每个文件逐个转换，转换完成后再转换下一个
 * 如果文件读取超时，返回失败
 * 如果文件读取成功，返回成功
 */

const successFileDataList = reactive([]); // 存储成功转换的文件数据
const handleUpload = (file_) => {
  onUploadLoading.value = true;
  return new Promise((resolve) => {
    let fileReadComplete = false; // 标志位，表示文件读取完成
    let fileErrMsg = ""; // 文件错误信息，如果正确则为空
    let successFileData = {};
    // ==================== 以下为判断文件读取是否完成的定时器 ====================
    // 设置一个超时定时器，以防文件读取失败
    const timer = setTimeout(() => {
      if (!fileReadComplete) {
        onUploadLoading.value = false;
        resolve({ status: "fail", error: "文件读取超时" });
      }
    }, 5000); // 假设 5 秒钟足够读取文件

    // 检查文件读取是否完成的循环;
    const checkCompletion = setInterval(() => {
      if (fileErrMsg !== "") {
        resolve({
          status: "fail",
          error: fileErrMsg,
        });
        onUploadLoading.value = false;
        clearInterval(checkCompletion); // 清除检查循环
      } else if (fileReadComplete) {
        resolve({
          status: "success",
          response: { url: window.location.href },
        });
        successFileDataList.push(successFileData);
        onUploadLoading.value = false;
        clearInterval(checkCompletion); // 清除检查循环
      }
    }, 1000); // 每 100 毫秒检查一次
    let file = file_[0];

    const reader = new FileReader();

    reader.onload = async (e) => {
      const data = new Uint8Array(e.target.result);
      const workbook = XLSX.read(data, { type: "array" });
      const sheetName = workbook.SheetNames[0];
      const worksheet = workbook.Sheets[sheetName];
      const json = XLSX.utils.sheet_to_json(worksheet, { header: 1 });
      // 延时1000ms模拟转换过程
      fileErrMsg = await checkFileValid(json);

      var result = await convertJsonToText(json);
      successFileData = {
        name: file.name,
        result: result,
        aesKey: "",
        ipfsPos: "",
      };
      console.log(successFileData);

      // 设置标志位为 true，表示文件读取完成
      fileReadComplete = true;
    };
    reader.readAsArrayBuffer(file.raw);
  });
};

const removeFile = (e) => {
  let file = e.file;

  // 从文件列表中删除该文件
  const index = files.value.findIndex((item) => item.name === file.name);
  if (index !== -1) {
    files.value.splice(index, 1);
  }
  // 从成功文件列表中删除该文件
  const successIndex = successFileDataList.findIndex(
    (item) => item.name === file.name
  );
  if (successIndex !== -1) {
    successFileDataList.splice(successIndex, 1);
  }
};

/**
 * 检查文件是否有效
 * @param json 读取的文件内容
 */
const checkFileValid = async (json) => {
  return new Promise((resolve) => {
    if (json.length === 0) {
      resolve("文件内容为空");
    } else if (json.length === 1) {
      resolve("文件只有一行");
    }
    if (json[0].length === 0) {
      resolve("列头为空");
    }
    // 检查是否有空值，且每一行的length应该相等
    for (let i = 0; i < json.length; i++) {
      if (i == 0) {
        // 检查数组是否有重复的列头
        let set = new Set(json[i]);
        if (set.size !== json[i].length) {
          resolve("列头名有重复");
        }
      }
      if (json[i].length !== json[0].length) {
        resolve(`第${i + 1}行的列数与第一行不相等`);
      }
      for (let j = 0; j < json[i].length; j++) {
        if (
          json[i][j] === undefined ||
          json[i][j] === null ||
          json[i][j] === ""
        ) {
          resolve(`第${i + 1}行第${j + 1}列单元格为空`);
        }
      }
    }
    resolve("");
  });
};

const checkFile = (file) => {
  if (file.size > 1024 * 1024 * 10) {
    return "文件大小不能超过10M";
  }
  return true;
};

const convertJsonToText = async (json) => {
  return json
    .map((row) =>
      row
        .map((cell) =>
          typeof cell === "string" ? cell.replace(/\s+/g, "_") : cell
        )
        .join(" ")
    )
    .join("\n");
};
// ==================== 下一步 ====================
const emit = defineEmits();
const handleNextStep = () => {
  // 将successFileDataList转为json字符串
  let jsonSuccessFileDataList = JSON.stringify(successFileDataList);
  console.log(jsonSuccessFileDataList);
  // 跳转到下一步
  emit("initFileOK", jsonSuccessFileDataList);
};

// ==================== 辅助函数 ====================
// 语言配置
const languageConfig = {
  upload: {
    progress: {
      uploadingText: "转换中",
      waitingText: "待转换",
      failText: "转换失败",
      successText: "转换成功",
    },
    triggerUploadText: {
      normal: "点击转换",
    },
  },
};
</script>

<style>
</style>