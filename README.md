# LINK Cinema sample service
 
LINK Cinema is a sample service of LINE Blockchain, demonstrating how to utilize LINE Blockchain Developers and BITMAX Wallet. It's an API server implemented using [Gin](https://github.com/gin-gonic/gin), a web framework written in Go.
 
Visit [LINE Blockchain Docs](https://docs.blockchain.line.me/sample-services/Link-cinema) for more details.
 
## Development environment
### Language
* Go 1.13+
### Back-end
* Gin
 
## Getting started
### Setting up the environment
You need Go 1.13 or higher to build the source code. Download and install the required version from [Go Downloads](https://golang.org/dl/).
 
LINK Cinema server brings the information of LINE Blockchain Developers from a separate configuration file written in TOML. The structure of the configuration object is as follows:
 
```
type APIConfig struct {
    LBDAPIEndpoint       string // Real API server address of LINE Blockchain Developers
    LINEAPIEndpoint      string // Real API server address of LINE
    LINEAccessEndpoint   string // Real API server address of LINE Access
    Endpoint             string // Service server address
    WalletAddress        string // Address of the service wallet
    WalletSecret         string // Secret key of the service wallet
    APIKey               string // API key of the service issued by LINE Blockchain Developers
    APISecret            string // API secret of the service issued by LINE Blockchain Developers
    ChannelID            string // ID of the channel issued by LINE Developers
    ChannelSecret        string // Secret key of the channel issued by LINE Developers
    ServiceContractID    string // Contract ID of the service token which is used as membership rewards points
    ItemContractID       string // Contract ID of the item tokens which are used as movie tickets or discount coupons
    FungibleTokenType    string // Token type of the non-fungible item tokens which are used as movie tickets or discount coupons
    NonFungibleTokenType string // Token type of the fungible item token which is used as movie tickets
    UserID               string // User ID of the service user's LINE account
}
```
 
LINK Cinema server reads the configuration file through the environment variable, `CONFIG_PATH`, during runtime. Designate the path of the configuration file with `CONFIG_PATH`.
 
```bash
$ export CONFIG_PATH={config toml file path}
```
 
### Building source code
 
```bash
$ make
```
 
### Running the LINK Cinema server
 
```bash
$ cinema
```
 
To check out the API endpoints provided by the LINK Cinema server, open the API reference file created by Swagger as follows:
 
```bash
$ {endpoint}/swagger/index.html
```
 
## How to contribute
 
See [CONTRIBUTING.md](CONTRIBUTING.md) for more details.
 
## License
 
```
Copyright 2020 LINE Corporation
 
LINE Corporation licenses this file to you under the Apache License,
version 2.0 (the "License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at:
 
  https://www.apache.org/licenses/LICENSE-2.0
 
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
License for the specific language governing permissions and limitations
under the License.
```
 
See [LICENSE](LICENSE) for more details.
