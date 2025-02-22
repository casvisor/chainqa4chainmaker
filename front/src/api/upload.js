import request from '@/api/request.js'
import { useApiStore } from '@/stores/api.js'
import { useUserStore } from '@/stores/user'
const baseUrl = '/upload'

/**
 * 从后端获取上传文件的AES密钥。（废弃：改为直接由前端生成）
 * @param {*} keyNum  需要生成的key数量
 * @returns  返回生成的key
 */
// export function getAesKey(keyNum) {
//     return request({
//         url: baseUrl + '/getAesKey',
//         method: 'post',
//         data: {
//             keyNum: keyNum
//         }
//     })
// }


export function uploadFile(aesKey, fileContent, fileName) {
    const apiStore = useApiStore()
    const userStore = useUserStore()
    return request({
        url: baseUrl + '/uploadFile',
        method: 'post',
        data: {
            uId: userStore.userId,
            aesKey: aesKey,
            fileContent: fileContent,
            fileName: fileName,
            apiUrl: {
                ipfsServiceUrl: apiStore.ipfsServiceUrl,
                chainServiceUrl: apiStore.chainServiceUrl,
                contractName: apiStore.contractName,
            }
        }
    })
}