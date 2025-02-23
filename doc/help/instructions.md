:::color3
由于腾讯云区块链设置了流量控制，1s 只能请求一次，请勿频繁点击

嵌入的网站地址请看第三部分

:::

<h1 id="OYV1e"> 使用说明</h1>
<h2 id="cDBlH">名词解释</h2>
+ 数据集：类似于 Mysql 的表/飞书的多维表格，其中第一行为列名（表头），有关数据集的通用规定，请看 2.1 节）
+ 分片：同一个数据集的不同切片。因为有些文件太大，所以可以进行按行拆分（比如 1~1000 行为一个分片，1001~1002 为一个分片）。注意，每一个分片的第一行需要为数据集表头，且各个分片的表头、格式应该一致。在本项目中，IPFS 存储的文件按”分片“为基本单位。
+ 日志：用户执行查询操作的日志

<h2 id="Aouzj">页面导引</h2>
<h3 id="anrOQ">首页</h3>
![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1740289568333-cacc55d9-86c2-4e95-b10e-37d2f0455449.png)

+ 主功能区：
    - 上传信息：上传分片以供查询
    - 查询信息：填写查询表单，完成查询
    - 下载文件：下载刚刚上传的分片信息，可用于二次验证是否上传成功。_<font style="color:#8A8F8D;">（此操作不经过区块链，也不经过查询模块，仅供验证是否上传成功。场景：老师查看成绩单是否上传成功，医生查看病例是否上传成功）</font>_
    - 审计日志：在查询信息模块的每一次查询都会记录到区块链上。

> 其中，上传信息和查询信息需要设置用户 ID（网页会校验）；下载文件和审计日志无需设置用户 ID
>

+ 其他功能区：
    - 点击”区块链复杂查询与审计模块“标题，可以返回主页
    - 点击用户名可修改![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1740289545448-302fac58-3130-4355-ab46-859a50f10c78.png)
    - 右上角![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736835073007-ad4dad38-1288-4225-a21d-5c6717cef1d4.png)分别是：帮助文档 、API 设置**<font style="color:#DF2A3F;">（请勿随意更改）</font>**、返回主页

<h3 id="ATtcJ">上传信息</h3>
<h4 id="T4qrq">第一步：上传</h4>
![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736835269872-9258e323-d268-4bf8-919a-202adc457e44.png)

点击左侧![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736835179764-e289ff7c-c6c2-4641-a4d3-bc499fd57010.png)按钮，选择分片文件（支持 Excel 格式）。系统会自动转换为适用于 IPFS 存储的文本格式。

若遇到格式错误，会进行提示（如上图第一行），格式正确会显示转换成功

上传成功后，点击![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736835280051-7dcd7ff2-dc2b-4b55-96ff-e3470e7c5ad1.png)进行转换确认

> 注意：
>
> + 上传的文件以分片（详见 1.1 名词解释）为基本单位。例如这里就有两份语文成绩单。是同一个语文成绩单的两个分片（详见 2.2.1 测试分片）。详细的分片格式要求，请参考 2.1 节
> + 不允许上传同名文件
> + 上传文件操作为浏览器本地操作
> + 格式错误的文件不会被计入后续的转换中。
>



<h4 id="Kt9wL">第二步：转换确认</h4>
点击下一步后，出现如图界面

![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736835476894-7fb0d154-22c8-462b-a45e-f325c87e2f97.png)

上面列举了每一份文件被转换后的效果。

转换的单元格之间以空格分开，行之间用换行符隔开。（<font style="color:rgb(38, 38, 38);">详细的分片被转换的格式说明，请参考 2.1 节</font>）

点击左侧三角![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736835569089-9dd42782-3956-4f36-a24e-de019300f072.png)可进行折叠；

点击右侧下载明文可以下载文本 txt，查看密钥可查看前端安全生成的 md5 密钥（用于后续加密）

若确认转换结果，请划到最下面，点击![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736835638179-c8abb0c1-8c86-40b2-aa51-30e8c662fc35.png)即进行上传操作；若发现异常，请点击上一步重新上传

<h4 id="AzwE1">第三步：上传</h4>
![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736836369575-945c8593-fba6-416d-9ff0-00ba2127e34e.png)

进入该页面后，会**自动执行**上传任务，等待结束即可。

由于腾讯云流控限制，上传任务每 3s 触发一次，所以需要等待。

在上一步中的每个文件分片将会被加密密钥进行 AES 加密后上传到 IPFS 系统中。系统返回 IPFS 的地址信息。

点击![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736836500644-cd3b17c0-0b74-41bd-8088-24dcc0e77f89.png)可下载上传文件的密文或明文，相关操作同 1.2.3.2 节。

点击![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736836538986-ab95680a-4fb4-4eb3-810d-876ffbbbb44a.png)可以复制相关信息到剪切板，例如：

```plain
文件名：7人测试数据（语文成绩单）第一份.xlsx
IPFS位置：QmcZsK2j9EDoqbYBgEZoLfNH6BeAoRz5rVTaYjVH5zm8Bu
加密密钥：386e5ca5263273da8f0bbd1df33bb10e
```

:::color3
请妥善保管相关信息，**<font style="color:#DF2A3F;">特别是IPFS位置</font>**。这是查询该数据的必备信息！本项目由于时间原因不存储该数据。

:::



<h3 id="vNLqy">下载文件</h3>
首页点击![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736836936026-1aa50532-e20a-4794-a4f9-0d273d408400.png)可用于验证是否上传成功。

<h4 id="PAWYs">填写下载表单</h4>
![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736836920525-682dfd46-bdc1-4cb3-8231-0a6f060e532a.png)

其中 IPFS 地址为必填，AES 密钥（就是 1.2.2.2 和 1.2.2.3 节所展示的加密密钥）选填

如果填写了 AES 密钥，则可以下载明文数据，否则只能下载密文数据

<h4 id="nGmgB">选择下载类型（明文/密文）</h4>
![传入AES密钥效果](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736837075239-9f4265d7-d0fc-4cc1-95cf-5cba9f15a26a.png)![未传入效果](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736837092738-15ec08de-bf7f-44a2-8123-4ffb25060f81.png)

![下载结果](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736837157192-0d675fe7-b09a-4eb0-8f14-a51f51bb3df2.png)

如果传入**<font style="color:#DF2A3F;">错误的密钥</font>**，点击下载密文时会进行提示报错，请关注页面上方，例如：![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736837187171-5586d091-94d1-4ec5-9b00-62dd6ce23e7b.png)

> 注意：
>
> + <font style="color:#000000;">此操作不经过区块链，也不经过查询模块，仅供验证是否上传成功。场景：老师查看成绩单是否上传成功，医生查看病例是否上传成功</font>
> + <font style="color:#000000;">下载密文后</font>**<font style="color:#DF2A3F;">不能自主解密（如通过其他网站的 AES 解密工具解密）</font>**<font style="color:#000000;">，</font>**<u><font style="color:#2F4BDA;">均需通过本接口/查询接口解密</font></u>**<font style="color:#000000;">，因为AES的密钥填充、IV值等情况略有不同</font>
>

<h3 id="Ohuos">查询信息</h3>
![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736837549434-4bf2aeb5-2757-4fd7-a8b4-80625637eb00.png)

<h4 id="kW3rR">单表查询</h4>
单表查询仅涉及到一个数据集的多个分片查询。不涉及数据集的联表操作。

同一个数据集的分片表头信息应该完全一致！

<h5 id="dweaL">联表类型</h5>
请选择”单表查询“

<h5 id="UtEnt">分片与展示</h5>
![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736837701259-91cc91a4-3aa9-4f15-ab79-fb750630c6b8.png)

因为只涉及到一个数据集，所以不能创建或者删除数据集

+ <font style="color:rgba(0, 0, 0, 0.9);">IPFS分片地址：请填写分片地址，按</font>**<font style="color:rgba(0, 0, 0, 0.9);">回车</font>**<font style="color:rgba(0, 0, 0, 0.9);">确认一个分片。为保证查询效率，最多 3 个分片。</font>

> 例如图中的两个 IPFS 地址都是语文成绩单的地址，分别是第一分片和第二个分片
>

+ <font style="color:rgba(0, 0, 0, 0.9);">数据集展示范围：可选择展示本数据集全部列/选择部分列展示/不展示列。此选择不会影响查询匹配的结果条目个数，仅会影响你看到的结果</font>
+ <font style="color:rgba(0, 0, 0, 0.9);">数据集展示列：仅”选择部分列展示“需要输入，按回车确认每个列。</font>**<font style="color:#ED740C;">若输入的列并不存在，那么将会忽略。</font>**

> 例如：你选择展示全部列，那么会返回语文成绩单的所有列（包括<font style="color:rgba(0, 0, 0, 0.9);"> id、name 和 Cscore 列</font>）；如果你选择<font style="color:rgba(0, 0, 0, 0.9);">部分列展示，且选择 name 和 Cscore 列，那么在查询结果（见 1.2.4.4 节）中将看不到 id 列的内容。</font>
>
> <font style="color:rgba(0, 0, 0, 0.9);">但是，你仍然可以选择 Id 列作为查询条件的列。（见 1.2.4.1.3 节）</font>
>

<h5 id="kwHIK">查询条件</h5>
![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736838131695-c2b41ae8-6b4b-47b1-8a18-bc16b417d310.png)

![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736838095787-f58840f5-2d05-4134-99e4-4e0d9666176a.png)

查询条件以”查询条件组“为基本单位。

+ **逻辑关系**：
    - 同一个条件组内的所有条件为 **AND** 关系
    - 条件组之间（如上图的条件组 1、条件组 2）为 **OR** 逻辑
    - 即：若某行数据完全满足**任一**条件组内的**所有**条件，该行数据匹配成功，流转至”分片与展示“<font style="color:rgba(0, 0, 0, 0.9);">（见 1.2.4.1.2 节）</font>部分进行展示范围的筛选。
+ **条件组编辑**
    - 默认为空条件组，展示为![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736838238320-17873316-efb1-490b-a81a-bd1e8a57df7b.png)，空条件组即标识”不设置匹配条件“，查询所有数据。
    - 你可以点击![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736838304290-7334c930-78c6-46d3-8696-e6cbed773b2d.png)来新增条件组，当已经存在至少一个条件组时，也可点击右侧的![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736838186299-21eb7407-2303-48fc-a2bb-a517639db367.png)进行新增。也可以点击条件组标签旁的![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736838203501-b53e04ef-d6fd-47eb-aaca-c176b6a3e085.png)删除条件组
+ **条件编辑**
    - ![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736838131695-c2b41ae8-6b4b-47b1-8a18-bc16b417d310.png?x-oss-process=image%2Fformat%2Cwebp)
    - 点击下方的![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736838544853-752e4aa1-d5d2-41b1-882e-007e18f2dbed.png)可在该组内新增一个条件；点击![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736838560430-c2da3a77-3149-41ac-950e-aa4602f6e53e.png)可删除该条件。若某个条件组**仅一个条件**，不允许删除该条件，若需删除，请点击![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736838592575-ef71a29d-7a49-4f52-8f55-46cf65c985f1.png)删除整个条件组。
    - 每个条件包括：
        * 列名：要查询的列名称，必须与分片表头的列名严格一致
            + 这里的列名可以手动输入，但是必须是分片所包含的列，否则会报错（见 1.2.4.4.2 节）
        * 列所在数据集：单表查询请忽略
        * 操作符和类型：请见 2.1.2 节
        * 基准值：要比较的值
    - 即<font style="color:rgb(38, 38, 38);">寻找</font>**<font style="color:rgb(237, 116, 12);">目标值（列名的每一行）</font>**<font style="color:rgb(237, 116, 12);">[ 操作符 ]  </font>**<font style="color:rgb(237, 116, 12);">基准值</font>**<font style="color:rgb(38, 38, 38);">的行。如上图就是查询</font>`<font style="color:rgb(38, 38, 38);"> Cscore </font>`<font style="color:rgb(38, 38, 38);">列大于 94 的行</font>

<h4 id="Akv5y">联表查询</h4>
<font style="color:rgb(38, 38, 38);">联表查询涉及到多个数据集的多个分片查询。暂时仅支持单个条件的联表，只支持内联 INNER 操作</font>

<font style="color:rgb(38, 38, 38);">同一个数据集的分片表头信息应该完全一致！</font>

<h5 id="g5Fzz"><font style="color:rgb(38, 38, 38);">联表类型</font></h5>
请选择“联表查询”

<h5 id="SN4S3">分片与展示</h5>
![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736839329670-c15cd851-230e-4803-8f49-94806ce957b5.png)

![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736839321655-148dabd8-37e2-4cd4-85a0-b2720941f61c.png)

由于联表，需要至少输入两个数据集。点击![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736839361321-c54627c0-887a-4c95-8038-73a26f182d8c.png)可新增数据集，点击![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736839369659-c44eb84a-b5de-416c-a994-aa87a36c2413.png)可删除数据集

同理可以配置各个数据集的展示范围。同 1.2.4.1.2 节

<h5 id="B1Jif">查询条件</h5>
基本同单表查询（见 1.2.4.2.2 节）

唯一不同的是需要选择列所在数据集

![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736839587126-77c61fc7-f116-4748-a5a4-7c62497ca5fe.png)

比如上面的意思就是查询：数据集 2 中 Mscore 列大于 82 **且** 数据集 1 中的 Cscore 列大于 92 的行集合。

> 若一个条件组中出现多个不同的数据集，则会**先联表后再筛选条件**。
>
> 比如上图中，若 Mscore<font style="color:rgb(38, 38, 38);">大于 82 但 Cscore=90，那么这两行即便联表成功，也不会被匹配。</font>
>

<h5 id="Rj2j8"><font style="color:rgb(38, 38, 38);">联表条件</font></h5>
<h6 id="pzvRO">联表条件解释</h6>
联表类似于 mysql 的 join

![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736839709983-92c9d8e5-7d27-45fc-9cf6-0a24d2328ebd.png)

联表条件告诉查询模块如何将两个不同的数据集合并为一个大的数据集。类似于查询条件

上图的意思是：数据集 1 的 id 和数据集 2 的 id 如果相同，那么就会被视为一行数据。

输入联表条件时，**<font style="color:#ED740C;">请遵从联表条件规范（下一节）</font>**，可以不用考虑顺序，程序会自动运行[拓扑排序](https://baike.baidu.com/item/%E6%8B%93%E6%89%91%E6%8E%92%E5%BA%8F/5223807)并进行校验

联表类型目前仅支持 Inner Join

<h6 id="QW6r7">联表条件规范</h6>
+ 所有数据集都应该被联表，不允许出现落单数据集
+ 两个数据集的联表条件唯一。不能出现环，也不能出现复合条件（比如输入了 A.id=B.id 为联表条件，又输入 A.name=B.name 为联表条件）

在以上两个规范下，`联表条件===数据集个数 -1`

当出现多个数据集（>=3）时，可能会由于错误输入出现联表环，会进行提示![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736839868332-dbadc1b0-effa-4330-850f-d152a52ad8b5.png)。例如：![这里数据集1和数据集2的联表条件出现了两次，出现环了](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736839897620-caf59305-4e4d-4e01-b461-6a6f02f23b22.png)



<h4 id="amyev">导入查询条件与导出查询条件</h4>
<h5 id="ejU92">导出查询条件</h5>
点击最下方生成 JSON 查询项可以导出查询条件（JSON 格式）到剪切板

![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736840369049-696e3b67-4566-4ef1-bae6-d63cb162479a.png)

（右边的生成字符串查询项按钮在最新版已经被删除）



<h5 id="vqKXT">导入查询条件</h5>
点击![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736842230184-e6a2aa73-581e-4ee6-8c2d-177c884753bc.png)，弹出窗口可输入查询条件 JSON

场景：老师分享查询条件时导入

![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736842238127-7bc1f711-8762-405f-b0df-322ac4c53341.png)

:::color3
查询条件 JSON 有严格的格式规定。建议是按照 1.2.4.3.1 节的方式导出后导入。

或者直接用界面进行可视化设置查询条件

:::

<h4 id="ulJ3z">查询结果</h4>
<font style="color:rgb(38, 38, 38);">由于腾讯云流控限制，分片下载任务每 2s 触发一次，需要等待。</font>

<h5 id="ZZjje">查询成功</h5>
![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736845374267-dd3ca7ea-2d6e-4056-b8cc-6f829180a09d.png)

查询成功会出现如上信息。

其中，若为联表查询，那么列名为`数据集分片地址_列名`以规避不同数据集存在相同列的情况，单表查询不会出现数据集名

可以点击![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736845448345-e3376991-3aa2-44ff-ba72-36e6dc184664.png)下载 excel

下方给出了查询结果 JSON、查询条件 JSON 等，可进行![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736845473703-a47fd8ec-3a34-40b7-80da-73d3e8540a4f.png)操作，点击![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736845479413-f0b87dad-49cf-466d-908b-541fcd986194.png)可展开

<h5 id="HbITE">查询出现错误</h5>
如果在条件组中输入不存在的列，比如输入“notExist”会报错。如下：

![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736845703974-cb180cc9-6ae5-435f-9410-11fc950e219c.png)

> <font style="color:rgb(38, 38, 38);">IPFS 文件不存在、IPFS 文件解析错误、列名不存在、类型不匹配（比如 string 类型列标明为 int）、正则表达式不正确、运算符不正确 </font>
>
> <font style="color:rgb(38, 38, 38);">发生上述错误时，会在查询结果的 message 中提到。其他错误请进行排查。</font>
>

<h3 id="fuwTS">审计日志</h3>
<h4 id="sCDkd">按用户 ID 查询</h4>
![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736845905214-2aee95e2-bfce-44c5-a043-d8e3bba190a6.png)

输入后点击![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736845920451-a3ccf684-9088-4fee-a69f-928fa0558ed5.png)即可查询，点击![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736845928231-979dd494-3804-4d33-897e-ead2ac987ff9.png)可下载表格数据

由于<font style="color:rgb(38, 38, 38);">查询条件和查询结果内容较多，以弹窗形式解析展现。</font>点击![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736845936251-4cc3574b-e3eb-4812-915c-4a14eef21d11.png)可查看本行日志的查询条件和查询结果，如下。

![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736845968392-9fe1820f-1faf-4fe1-9c4d-3b9adc67da9d.png)

![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736846014084-19db15c6-0a3b-4e36-9185-662656a3c08b.png)

点击弹窗右上角的![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736846045461-57f5fd86-5bcf-4770-84f6-c4c76a449220.png)可切换展示



查询失败的日志也会记录，如：

![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736846150695-fcb496c9-0fdc-40a4-88fa-f4a312d62e3f.png)

<h4 id="xcaA6">按时间范围查</h4>
同理：

![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736846121785-2ac59939-5869-4898-ae86-75cb8948b188.png)

> 提示：请尽量查询 1 月 12 日及之后的数据，之前的数据由于不适配等信息，在“查看”弹窗中可能无法展示。
>

<h1 id="C3Gsz">规范与演示</h1>
<h2 id="NqIeT">规范</h2>
<h3 id="llUsd">数据集/分片规范</h3>
上传的文件需满足如下规范：

+ 第一行必须是列名，列名不能仅为`*`
+ 单元格不允许为空，值不能包含空格，
+ 同一列的类型应相同（若不能保证相同，请在查询时使用 string 类型查询，此时比较大小将按字典序）
+ 有几列，那么每行数据就有几个单元格（严格相等，不能多不能少）

> 可以仿照 Mysql/飞书多维表格/Notion 的 database 来看待
>

<h3 id="LOrIU">可选操作与类型</h3>
<h4 id="VroL8">操作符</h4>
+ <font style="color:rgb(38, 38, 38);">eq：相等	</font>
+ <font style="color:rgb(38, 38, 38);">ne：不等</font>
+ <font style="color:rgb(38, 38, 38);">gt：> （string 类型将按字典序进行比较）——即寻找</font>**<font style="color:rgb(237, 116, 12);">目标值</font>**<font style="color:rgb(237, 116, 12);">大于</font>**<font style="color:rgb(237, 116, 12);">基准值</font>**<font style="color:rgb(38, 38, 38);">的行</font>
+ <font style="color:rgb(38, 38, 38);">lt：< （string 类型将按字典序进行比较）</font>
+ <font style="color:rgb(38, 38, 38);">ge：>=</font>
+ <font style="color:rgb(38, 38, 38);">le：<=</font>
+ <font style="color:rgb(38, 38, 38);">regexp：正则表达式匹配（程序会判断正则是否正确）</font>
+ <font style="color:rgb(38, 38, 38);">contain：是否包含某个子串</font>
+ <font style="color:rgb(38, 38, 38);">suffix：后缀</font>
+ <font style="color:rgb(26, 32, 41);">prefix：前缀</font>

> 注：<font style="color:rgb(38, 38, 38);">regexp、contain、suffix、</font><font style="color:rgb(26, 32, 41);">prefix 四个操作仅 string 类型可用，其余公用</font>
>

<h4 id="JaLlP">比较类型</h4>
+ <font style="color:rgb(38, 38, 38);">string（按字典序比较）</font>
+ <font style="color:rgb(38, 38, 38);">int（按整数值比较）</font>
+ <font style="color:rgb(38, 38, 38);">float（按浮点数进行比较）</font>





<h2 id="anPR9">测试数据</h2>
<h3 id="o8luD">测试分片</h3>
[7人测试数据（语文成绩单）第一份.xlsx](https://www.yuque.com/attachments/yuque/0/2025/xlsx/40383230/1736836237373-61482636-6803-455e-ba62-879a0c43c8bd.xlsx)

[7人测试数据（语文成绩单）第二份.xlsx](https://www.yuque.com/attachments/yuque/0/2025/xlsx/40383230/1736836237257-1899c4d6-f3a3-44da-9c04-7a5a23871531.xlsx)

[7人测试数据（数学成绩单）全.xlsx](https://www.yuque.com/attachments/yuque/0/2025/xlsx/40383230/1736836237374-15646969-36d4-4bd2-a69c-28d6c6d62299.xlsx)

[7人测试数据（错误成绩单）.xlsx](https://www.yuque.com/attachments/yuque/0/2025/xlsx/40383230/1736836251618-fb9f8726-8561-45c5-9f16-53676d764598.xlsx)



<font style="color:rgb(51, 51, 51);">文件名：7人测试数据（语文成绩单）第一份.xlsx</font>

<font style="color:rgb(51, 51, 51);">IPFS位置：QmWCKFXeBmWfKq9MHjiP95miz3f8k27Ub986gfnTMivyau</font>

<font style="color:rgb(51, 51, 51);">加密密钥：cd88880be5f1d08d4f26e3429733c222</font>

<font style="color:rgb(51, 51, 51);"></font>

<font style="color:rgb(51, 51, 51);">文件名：7人测试数据（语文成绩单）第二份.xlsx</font>

<font style="color:rgb(51, 51, 51);">IPFS位置：QmfAuB5ed5QCUcFs3tBUiNbxPRpsbQSje9QzUQNhVyhE64</font>

<font style="color:rgb(51, 51, 51);">加密密钥：94784b3ddd9ea291ab4ea80c75d0c773</font>

<font style="color:rgb(51, 51, 51);"></font>

<font style="color:rgb(51, 51, 51);">文件名：7人测试数据（数学成绩单）全.xlsx</font>

<font style="color:rgb(51, 51, 51);">IPFS位置：QmNuoRoAKHLZBuyyCBmiajKn4YQJiWpvGMtiFtFdonEAHH</font>

<font style="color:rgb(51, 51, 51);">加密密钥：da673d2b07eca446690cfdfccd7ee042</font>

<h3 id="DoSJz">示例查询条件 JSON</h3>
```json
{"queryConcatType":"multi","filePos":[["QmQmELkad812TKSyxM63u7Z5ymPMup452HoruejBRvNodF","QmZoYQvTVtTBaVcgGankWfeAAwAD9QaC4G8y1V4ZKb1ywz"],["QmSWK127q4EsVw5qhZ3MmdNiUQPmL158fLFXsHy9x3Ft58"]],"returnField":["QmQmELkad812TKSyxM63u7Z5ymPMup452HoruejBRvNodF_Cscore","QmQmELkad812TKSyxM63u7Z5ymPMup452HoruejBRvNodF_id","QmSWK127q4EsVw5qhZ3MmdNiUQPmL158fLFXsHy9x3Ft58_Mscore"],"queryConditions":[[{"field":"Cscore","pos":"QmQmELkad812TKSyxM63u7Z5ymPMup452HoruejBRvNodF","compare":"gt","val":"92","type":"int"},{"field":"Mscore","pos":"QmSWK127q4EsVw5qhZ3MmdNiUQPmL158fLFXsHy9x3Ft58","compare":"gt","val":"86","type":"int"}],[{"field":"id","pos":"QmQmELkad812TKSyxM63u7Z5ymPMup452HoruejBRvNodF","compare":"eq","val":"1","type":"int"}]],"jointConditions":[{"pos1":"QmQmELkad812TKSyxM63u7Z5ymPMup452HoruejBRvNodF","field1":"id","pos2":"QmSWK127q4EsVw5qhZ3MmdNiUQPmL158fLFXsHy9x3Ft58","field2":"id","compare":"eq","type":"int","jointType":"INNER"}]}
```

<h3 id="QypRH">演示视频</h3>
[演示.mp4](https://www.yuque.com/attachments/yuque/0/2025/mp4/40383230/1736847724604-3053106d-b2a6-4d40-aefa-b6a7f608e572.mp4)

<h1 id="rydIZ">网站地址</h1>
网页嵌入：

直接跳转到 [http://47.113.204.64:9000/](http://47.113.204.64:9000/) 即可（跳转页面为首页）

+ 若需传入 userId，请使用[http://47.113.204.64:9000/?userId=TEST](http://47.113.204.64:9000/?userId=TEST)
    - TEST 即用户 id，若不传，会强制弹窗要求传入



若想跳转到具体的功能页，请使用如下网址

+ 查询模块：[http://47.113.204.64:9000/query?userId=TEST](http://47.113.204.64:9000/query?userId=TEST)
    - TEST 即用户 id，若不传，会强制弹窗要求传入![](https://cdn.nlark.com/yuque/0/2025/png/40383230/1736848152537-9288518e-2fff-4d23-b5f1-aed7b47243b5.png)
+ 上传模块：[http://47.113.204.64:9000/upload?userId=TEST](http://47.113.204.64:9000/upload?userId=upload_TEST)
    - <font style="color:rgb(38, 38, 38);">TEST 即用户 id，若不传，会强制弹窗要求传入</font>
+ <font style="color:rgb(38, 38, 38);">下载文件：</font>[http://47.113.204.64:9000/download](http://47.113.204.64:9000/download) 该模块无需 userID
    - 当然也可以使用[http://47.113.204.64:9000/download?userId=TEST](http://47.113.204.64:9000/download?userId=TEST)来传入 ID
+ 审计日志：[http://47.113.204.64:9000/log](http://47.113.204.64:9000/log) <font style="color:rgb(38, 38, 38);">该模块无需 userID</font>
    - 当然也可以使用[http://47.113.204.64:9000/log?userId=TEST](http://47.113.204.64:9000/log?userId=TEST)来传入 ID

