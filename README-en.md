# ChainQA4chainmaker

[简体中文](./README.md) | English

ChainQA (Blockchain Audit and Query Platform) for Tencent Cloud - [ChainMaker](https://chainmaker.org.cn/home) (Version 2.2)

> This document is partially translated by a large language model. The documents linked in this article only support Chinese. Thank you for your understanding.

## Project Introduction

ChainQA is a component based on blockchain and IPFS technologies, integrating data encryption storage, efficient query, and audit functions. It aims to solve the problems of data security and privacy protection.

Users only need to upload table files that meet the established format requirements. ChainQA will automatically encrypt them and store them in IPFS, and at the same time, store the corresponding keys securely on the blockchain. For data query, when the uploader shares the IPFS address, other users can initiate a query operation. At this time, the decryption process is completed inside the blockchain, and the query results are returned without leaking the encryption key throughout the process. All query operations will be recorded, and audit tracing can be performed at any time according to the query account or time range.

ChainQA supports multi-condition and joint table query functions, allowing users to obtain the required data flexibly and accurately. This version only requires the querier to obtain the IPFS address for querying. Secondary development can be carried out independently to add auxiliary functions such as query passwords.

## Project Architecture

### Main Project Architecture

- back: Backend code, mainly responsible for basic operations such as IPFS encryption and decryption, converting Excel data, and encapsulating requests.
- front: Frontend code for page display.
- contract: Contract code deployed on [Tencent Cloud ChainMaker](https://console.cloud.tencent.com/tbaas/chainmaker/chain/chainmaker-demo/basicInfo?chainId=1&demo=1). **Please note that this contract code is only compatible with ChainMaker versions 2.1 to 2.2.**
- doc: Some documents (Excel) that can be used for testing, some pictures, and some help documents.
- tencent-chainmaker: SDK for connecting the backend with Tencent Cloud Blockchain. **Configuration is required.**
- docker-compose.yml: Docker orchestration file for starting Docker.

### Backend (Back)

👋 If you just want to **experience the backend**, please run it as follows:

Install the dependency packages:

```bash
go mod vendor
```

Run:

```bash
go run main.go
```

The backend listens on port 9000.

👋 If you want to deploy the **entire project**, please refer to the deployment section.

### Frontend (front)

👋 If you just want to **experience the frontend**, please run it as follows:

You need to install node and yarn in advance.

Run (in the frontend debugging environment):

```bash
yarn dev
```

👋 If you want to deploy the **entire project**, please refer to the deployment section.

### Blockchain Smart Contract (contract)

It is located in the contract folder.

Among them, `tencentChainqaContractV221demo01.7z` is the compiled version. `src` is the source code folder.

> Please note that the contract code in this file is only compatible with ChainMaker versions 2.1 to 2.2 and can be deployed to [Tencent Cloud ChainMaker](https://console.cloud.tencent.com/tbaas/chainmaker/chain/chainmaker-demo/basicInfo?chainId=1&demo=1).

The recommended contract name is: `tencentChainqaContractV221demo01` (otherwise, you need to change the contract name in the frontend).

![Contract Methods](./doc/img/contract-method.png)

### tencent-chainmaker

It is mainly a system for the backend to **interact (call contracts)** with [Tencent Cloud ChainMaker](https://console.cloud.tencent.com/tbaas/chainmaker/chain/chainmaker-demo/basicInfo?chainId=1&demo=1). It is a secondary encapsulation of the blockchain interaction SDK provided by Tencent Cloud.

This part needs to be configured:

1. Go to the [Tencent Cloud API](https://console.cloud.tencent.com/cam/capi) page, apply for a key, and save the SecretId and SecretKey.
2. Enter the `/tencent-chainmaker/conf/config_template.ini` file, fill in the corresponding parts, and then copy a copy and save it as `config.ini`.

> The key is a high-risk part! It cannot be disclosed in any form! It needs to be properly kept.

> Why is this part needed? Because the overall link is as follows:

> Frontend —— Backend —— tencent-chainmaker —— ChainMaker

## Project Docker Environment Configuration and Deployment Guide

This project has built-in Docker-related configuration files and supports running the project in a Docker environment. The following will elaborate on how to configure Docker and complete the packaging and running operation process.

### Step 1: Package the Frontend (Optional)

If you have modified the frontend source code, you need to repackage the frontend project according to this step. The specific operations are as follows:

1. Install frontend tools: You need to install two frontend tools, [node](https://nodejs.org/zh-cn) and [yarn](https://yarnpkg.com/).
2. Switch directories and execute the packaging command: Switch to the frontend directory (`/front`) and run the `yarn build` command.
3. Generate the packaged files: Wait a moment, and the frontend packaged files will be generated in the `/front/dist` folder. Note that the frontend packaged files are already stored in the original `/front/dist` directory of the project files. If you have not modified the source code, you can use them directly; **if there are modifications, you need to package them**.

### Step 2: Move the Frontend Packaged Files to the Backend Folder

Copy the entire `/front/dist` directory to the `/back/dist` directory. After completing this operation, no additional configuration is required.

### Step 3: Configure tencent - chainmaker

Please complete the configuration of SecretId and SecretKey according to the instructions in the tencent - chainmaker part of this document.

### Step 4: Install the Blockchain Smart Contract

Configure according to the guidelines in the blockchain smart contract part of this document.

### Step 5: Start the Container

Switch to the project root directory and run the following command to start the container:

```bash
docker-compose up -d
```

This command will automatically start the `docker-compose.yml` file in the root directory. This `yml` file defines two containers, namely the backend container and the tencent - chainmaker container. Since the frontend has been packaged and placed in the backend in the second step, no additional frontend container is required. After that, the system will load the images and start the containers according to `back/Dockerfile` and `tencent-chainmaker/Dockerfile` respectively.

If you have not changed the default configuration, the backend service will run on port 9000, and the tencent - chainmaker service will run on port 9001.

### Step 6: Run the Project

Open a browser and enter `localhost:9000` in the address bar to access the project page.

**Note**: If you need to modify the backend port number, please edit the `port` configuration item in the `/back/conf/conf.ini` file.

### Step 7: Modify the API Address

Your API address may be different from mine, so you need to modify it.

![img](./doc/img/modify-in-front.png)

You need to configure and modify three addresses and the contract name:

1. **Backend Container API Address**: Only change the "Host IP:Port Number" (default is 47.113.204.64:9000). Use the server IP or local "localhost" for the host IP. If you have not modified the `/back/conf/conf.ini` file, use port 9000.
2. **IPFS Container Address**: This document does not cover the installation and configuration of the IPFS container. You can refer to `/doc/help/Deploying-IPFS.md` or [related documents](https://www.yuque.com/jjq0425/pku/cm112pwu470v3q9n). After installation, set the IP.
3. **tencent - chainmaker Container API Address**: Only change the "Host IP:Port Number" (default is 47.113.204.64:9000). On a local computer, try `host.docker.internal`, `172.17.0.1`, or `localhost` in sequence.
4. **Contract Name**: Use the contract name filled in Tencent Cloud Blockchain.

> Note:
>
> After completing the above configuration modifications and clicking save, the configuration you made is **only valid in the current browsing session**. Once you refresh the browser, the above configuration will return to the default state, and you need to reconfigure it.
>
> If you want the frontend to automatically apply the modified configuration and avoid the cumbersome operation of manual configuration every time, you can make corresponding adjustments to the frontend's `\front\src\stores\api.js` file after completing the configuration modification. Then, re-execute the deployment process (including frontend packaging). In this way, the frontend can automatically load and use the modified configuration to achieve persistent application of the configuration.

## Project Usage Instructions

Please refer to [this article](https://www.yuque.com/jjq0425/pku/acgk6l3ax4gf98wh) or view `/doc/help/instructions.md`.
