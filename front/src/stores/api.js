import { ref } from 'vue';
import { defineStore } from 'pinia';

export const useApiStore = defineStore('api', () => {
    // const api = ref('http://localhost:9000/api');
    const api = ref('http://47.113.204.64:9000/api');
    const ipfsServiceUrl = ref('http://47.113.204.64:5001');
    // const chainServiceUrl = ref('http://localhost:8088/chain/contract_invoke');
    // const chainServiceUrl = ref('http://host.docker.internal:9001/tencent-chainapi/exec');
    const chainServiceUrl = ref('http://47.113.204.64:9001/tencent-chainapi/exec');
    // const contractName = ref('chainQA_test_offchain_0_001');
    const contractName = ref('tencentChainqaContractV221demo01');


    function setApi(url) {
        api.value = url;
    }

    function setIpfsServiceUrl(url) {
        ipfsServiceUrl.value = url;
    }

    function setChainServiceUrl(url) {
        chainServiceUrl.value = url;
    }

    function setContractName(name) {
        contractName.value = name;
    }

    return { api, setApi, contractName, chainServiceUrl, ipfsServiceUrl, setIpfsServiceUrl, setContractName, setChainServiceUrl };
});
