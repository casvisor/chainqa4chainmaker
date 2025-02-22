import request from '@/api/request.js'
import { useApiStore } from '@/stores/api.js'
import { useUserStore } from '@/stores/user.js'
const baseUrl = '/query'



export function queryData(qItem, fPoses) {
    const apiStore = useApiStore()
    const userStore = useUserStore()
    return request({
        timeout: 240000,//240s超时
        url: baseUrl + '/queryData',
        method: 'post',
        data: {
            uId: userStore.userId,
            queryItem: qItem,
            // filePoses: fPoses,
            apiUrl: {
                ipfsServiceUrl: apiStore.ipfsServiceUrl,
                chainServiceUrl: apiStore.chainServiceUrl,
                contractName: apiStore.contractName,
            }
        }
    })

}