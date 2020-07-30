# LINK Cinema

Link cinema is a sample DAPP to provide trials of LBD(LINE blockchain developers) and LBW(LINE blockchain wallet). It is an API server consisting of Gin (https://github.com/gin-gonic/gin), a golang webframework

## Getting Start

Requires [Go 1.13+](https://golang.org/dl/)

Building the source

```
    $ make
```

Link cinema must load the configuration feature from an external Toml file. The following code is the structure of configuration feature.

```
    type APIConfig struct {
        LBDAPIEndpoint       string //API endpoint of LBD
        LINEAPIEndpoint      string //API endpoint of LINE
        LINEAccessEndpoint   string //API endpoint of Line-Access
        Endpoint             string //Endpoint of this server
        WalletAddress        string //Address of service wallet
        WalletSecret         string //Secret key of service wallet
        APIKey               string //API key of LBD
        APISecret            string //API secret key of LBD
        ChannelID            string //LINE Developers channel ID
        ChannelSecret        string //Secret key of LINE developers channel
        ServiceContractID    string //Contract ID of service token using as movie token
        ItemContractID       string //contract ID of item token using as movie-ticket, movie-discount token
        FungibleTokenType    string //Token type of fungible token using as movie-discount token 
        NonFungibleTokenType string //Token type of non-fungible token using as movie-ticket token
        UserID               string //LINE user ID
    }
```
Link Cinema loads the toml file path set in environment variable `CONFIG_PATH` at runtime

```
    $ export CONFIG_PATH={config toml file path}
```

Run a Link Cinema

```
    $ cinema
```

Access to swagger

```
    $ {endpoint}/swwager/index.html
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


