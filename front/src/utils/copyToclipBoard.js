import { MessagePlugin } from "tdesign-vue-next"
const copyToClipBoard = async (text) => {
    let clipboard = {
        writeText: (text) => {
            let copyInput = document.createElement('input');
            copyInput.value = text;
            document.body.appendChild(copyInput);
            copyInput.select();
            document.execCommand('copy');
            document.body.removeChild(copyInput);
        }
    }
    if (clipboard) {
        await clipboard.writeText(text);
        MessagePlugin.success('复制成功');
    }

}

export default copyToClipBoard