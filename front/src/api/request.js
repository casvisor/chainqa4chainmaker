import axios from 'axios'
import { NotifyPlugin } from 'tdesign-vue-next';
import { useApiStore } from "../stores/api";



const request = axios.create({
    timeout: 25000 // 请求超时时间
})

request.interceptors.request.use(
    config => {

        // 在服务器测试时候不带baseUrl，在NetworkSetting已经带了
        if (config.NetworkSetting == null || !config.NetworkSetting) {
            const apiStore = useApiStore()
            config.url = apiStore.api + config.url
        }

        return config
    },
    error => {
        return Promise.reject(error)
    }
)
request.interceptors.response.use(
    response => {
        return response.data
    },
    error => {
        if (error.response?.status !== 200) {
            if (error.response?.data?.code != 0 && error.response?.data?.msg !== undefined) {
                NotifyPlugin.error(error.response.data.msg)
            }
        }
        return Promise.reject(error)
    }
)



export default request