import request from '@/api/request.js'
import { useApiStore } from '@/stores/api.js'
const baseUrl = '/log'




export function logByUid(uId) {
    const apiStore = useApiStore()
    return request({
        url: baseUrl + '/logByUid',
        method: 'post',
        data: {
            uId: uId,
            apiUrl: {
                ipfsServiceUrl: apiStore.ipfsServiceUrl,
                chainServiceUrl: apiStore.chainServiceUrl,
                contractName: apiStore.contractName,
            }
        }
    })

}

/**
 * 
 * @param {*} startTime  传入的是秒级数字时间戳
 * @param {*} endTime  传入的是秒级数字时间戳
 * @returns 
 */
export function logByTimeRange(startTime, endTime) {
    const apiStore = useApiStore()
    return request({
        url: baseUrl + '/logByTimeRange',
        method: 'post',
        data: {
            startTime: startTime,
            endTime: endTime,
            apiUrl: {
                ipfsServiceUrl: apiStore.ipfsServiceUrl,
                chainServiceUrl: apiStore.chainServiceUrl,
                contractName: apiStore.contractName,
            }
        }
    })

}