import request from '@/api/request.js'
import { useApiStore } from '@/stores/api.js'
const baseUrl = '/download'




export function downloadIPFSFile(cid) {
    const apiStore = useApiStore()
    return request({
        url: baseUrl + '/downloadIPFSFile',
        method: 'post',
        data: {
            cid: cid,
            apiUrl: {
                ipfsServiceUrl: apiStore.ipfsServiceUrl,
                chainServiceUrl: apiStore.chainServiceUrl,
                contractName: apiStore.contractName,
            }
        }
    })

}

export function tryDecryptFile(cid, aesKey) {
    const apiStore = useApiStore()
    return request({
        url: baseUrl + '/tryDecryptFile',
        method: 'post',
        data: {

            cid: cid,
            aesKey: aesKey,
            apiUrl: {
                ipfsServiceUrl: apiStore.ipfsServiceUrl,
                chainServiceUrl: apiStore.chainServiceUrl,
                contractName: apiStore.contractName,
            }
        }
    })

}