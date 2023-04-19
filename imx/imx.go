package imx

import (
	"context"
	"encoding/json"
	"github.com/immutable/imx-core-sdk-golang/imx"
	"github.com/immutable/imx-core-sdk-golang/imx/api"
	"github.com/immutable/imx-core-sdk-golang/imx/signers/ethereum"
	"github.com/immutable/imx-core-sdk-golang/imx/signers/stark"
	"log"
	"math/big"
	"nft/models"
	"strconv"
)

type Client interface {
	Close()
	CreateUser(ctx context.Context, user *models.User) (string, error)
	CreateCollection(ctx context.Context, info *CollectionInformation) error
	CreateMetadata(ctx context.Context, info *MetadataInformation) error
	CreateToken(ctx context.Context, info *MintInformation) error
	TransferToken(ctx context.Context, info *TransferInformation) error
	CreateOrder(ctx context.Context, info *OrderInformation) error
}

type IMX struct {
	client    *imx.Client
	l1signer  imx.L1Signer
	l2signer  imx.L2Signer
	chainId   *big.Int
	projectID int32
}

type ProjectInformation struct {
	ProjectName  string
	CompanyName  string
	ContactEmail string
}

type CollectionInformation struct {
	ProjectID       int32
	ContractAddress string
	CollectionName  string
	PublicKey       string
	MetadataUrl     string
}

type MetadataInformation struct {
	ContractAddress string
	Fields          []MetadataFieldInformation
}

type MetadataFieldInformation struct {
	Name string
	Type string
}

type MintInformation struct {
	ContractAddress string
	TokenID         string
	Blueprint       string
}

type OrderInformation struct {
	ContractAddress string
	TokenID         string
	Amount          uint64
}

type CreateDeposit struct {
	DepositAmountWei string
}

type TransferInformation struct {
	TokenID         string
	ContractAddress string
	ReceiverAddress string
}

func NewIMX(alchemyAPIKey string, l1SignerPrivateKey string, starkPrivateKey string, projectID int32) Client {
	apiConfiguration := api.NewConfiguration()
	cfg := imx.Config{
		APIConfig:     apiConfiguration,
		AlchemyAPIKey: alchemyAPIKey,
		Environment:   imx.Sandbox,
	}
	client, err := imx.NewClient(&cfg)
	if err != nil {
		log.Panicf("error in NewIMX: %v\n", err)
	}

	l1signer, err := ethereum.NewSigner(l1SignerPrivateKey, cfg.ChainID)
	if err != nil {
		log.Panicf("error in creating L1Signer: %v\n", err)
	}

	l2signer, _ := newStarkSigner(starkPrivateKey)

	return &IMX{client, l1signer, l2signer, cfg.ChainID, projectID}
}

func (i *IMX) CreateUser(ctx context.Context, user *models.User) (string, error) {
	l1signer, err := ethereum.NewSigner(user.Private, i.chainId)
	if err != nil {
		return "", err
	}

	l2signer, starkKey := newStarkSigner("")

	response, err := i.client.RegisterOffchain(ctx, l1signer, l2signer, user.Mail)
	if err != nil {
		return "", err
	}

	val, err := prettyStruct(response)
	if err != nil {
		return "", err
	}
	log.Println("RegisterOffchain response: ", val)

	// Get the accounts registered on offchain.
	usersResponse, err := i.client.GetUsers(ctx, l1signer.GetAddress())
	if err != nil {
		return "", err
	}
	log.Println("Registered accounts: ", usersResponse.GetAccounts())
	return starkKey, nil
}

func (i *IMX) Close() {
	i.client.EthClient.Close()
}

func newStarkSigner(privateStarkKeyStr string) (imx.L2Signer, string) {
	var err error
	if privateStarkKeyStr == "" {
		privateStarkKeyStr, err = stark.GenerateKey()
		log.Println("Stark Private key: ", privateStarkKeyStr)
		if err != nil {
			log.Panicf("error in Generating Stark Private Key: %v\n", err)
		}
	}

	l2signer, err := stark.NewSigner(privateStarkKeyStr)
	if err != nil {
		log.Panicf("error in creating StarkSigner: %v\n", err)
	}
	return l2signer, privateStarkKeyStr
}

func prettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func (i *IMX) CreateProject(ctx context.Context, info *ProjectInformation) (int32, error) {
	response, err := i.client.CreateProject(ctx, i.l1signer, info.ProjectName, info.CompanyName, info.ContactEmail)
	if err != nil {
		return -1, err
	}

	val, err := prettyStruct(response)
	if err != nil {
		return -1, err
	}

	log.Println("Created new project, response: ", val)

	// Get the project details we just created.
	projectResponse, err := i.client.GetProject(ctx, i.l1signer, strconv.FormatInt(int64(response.Id), 10))
	if err != nil {
		return -1, err
	}
	val, err = prettyStruct(projectResponse)
	if err != nil {
		return -1, err
	}

	log.Println("Project details: ", val)
	return response.Id, nil
}

func (i *IMX) CreateCollection(ctx context.Context, info *CollectionInformation) error {
	info.ProjectID = i.projectID
	info.PublicKey = i.l1signer.GetPublicKey()

	createCollectionRequest := api.NewCreateCollectionRequest(info.ContractAddress,
		info.CollectionName,
		info.PublicKey,
		info.ProjectID)

	createCollectionRequest.MetadataApiUrl = &info.MetadataUrl

	response, err := i.client.CreateCollection(ctx, i.l1signer, createCollectionRequest)
	if err != nil {
		return err
	}

	val, err := prettyStruct(response)
	if err != nil {
		return err
	}
	log.Println("Created new collection, response: ", val)

	// Get the collection details we just created.
	collectionReponse, err := i.client.GetCollection(ctx, info.ContractAddress)
	if err != nil {
		return err
	}
	log.Println("Created Collection Name: ", collectionReponse.Name)

	return nil
}

func (i *IMX) CreateMetadata(ctx context.Context, info *MetadataInformation) error {
	metadata := make([]api.MetadataSchemaRequest, 0)
	for _, f := range info.Fields {
		field := api.NewMetadataSchemaRequest(f.Name)
		field.SetType(f.Type)
		metadata = append(metadata, *field)
	}

	request := api.NewAddMetadataSchemaToCollectionRequest(metadata)
	response, err := i.client.AddMetadataSchemaToCollection(ctx, i.l1signer, info.ContractAddress, *request)
	if err != nil {
		return err
	}

	val, err := prettyStruct(response)
	if err != nil {
		return err
	}
	log.Println("Created new metadata, response: ", val)
	return nil
}

func (i *IMX) CreateToken(ctx context.Context, info *MintInformation) error {
	tokenID := info.TokenID
	tokenAddress := info.ContractAddress
	ethAddress := i.l1signer.GetAddress()
	blueprint := info.Blueprint

	var royaltyPercentage float32 = 1

	var mintableToken = imx.UnsignedMintRequest{
		ContractAddress: tokenAddress,
		Royalties: []imx.MintFee{
			{
				Percentage: royaltyPercentage,
				Recipient:  ethAddress,
			},
		},
		Users: []imx.User{
			{
				User: ethAddress,
				Tokens: []imx.MintableTokenData{
					{
						ID: tokenID,
						Royalties: []imx.MintFee{
							{
								Percentage: royaltyPercentage,
								Recipient:  ethAddress,
							},
						},
						Blueprint: blueprint,
					},
				},
			},
		},
	}

	request := make([]imx.UnsignedMintRequest, 1)
	request[0] = mintableToken

	mintTokensResponse, err := i.client.Mint(ctx, i.l1signer, request)
	if err != nil {
		return err
	}

	log.Printf("Mint Tokens response:\n%v\n", mintTokensResponse.Results[0].TxId)
	return nil
}

func (i *IMX) TransferToken(ctx context.Context, info *TransferInformation) error {
	signableToken1 := imx.SignableERC721Token(info.TokenID, info.ContractAddress)

	transferRequest1 := api.SignableTransferDetails{
		Amount:   "1",
		Receiver: info.ReceiverAddress,
		Token:    signableToken1,
	}

	batchTransferRequest := api.GetSignableTransferRequest{
		SenderEtherKey: i.l1signer.GetAddress(),
		SignableRequests: []api.SignableTransferDetails{
			transferRequest1,
		},
	}

	response, err := i.client.BatchNftTransfer(ctx, i.l1signer, i.l2signer, batchTransferRequest)
	if err != nil {
		return err
	}

	val, err := prettyStruct(response)
	if err != nil {
		return err
	}
	log.Println("Created new metadata, response: ", val)
	return nil
}

func (i *IMX) CreateOrder(ctx context.Context, info *OrderInformation) error {
	ethAddress := i.l1signer.GetAddress()                                    // Address of the user listing for sale.
	sellToken := imx.SignableERC721Token(info.TokenID, info.ContractAddress) // NFT Token
	buyToken := imx.SignableETHToken()                                       // The listed asset can be bought with Ethereum
	createOrderRequest := &api.GetSignableOrderRequest{
		AmountBuy:  strconv.FormatUint(info.Amount, 10),
		AmountSell: "1",
		Fees:       nil,
		TokenBuy:   buyToken,
		TokenSell:  sellToken,
		User:       ethAddress,
	}
	createOrderRequest.SetExpirationTimestamp(0)

	// Create order will list the given asset for sale.
	createOrderResponse, err := i.client.CreateOrder(ctx, i.l1signer, i.l2signer, createOrderRequest)
	if err != nil {
		return err
	}

	createOrderResponseStr, err := prettyStruct(createOrderResponse)
	if err != nil {
		return err
	}
	log.Printf("CreateOrder response:\n%v\n", createOrderResponseStr)
	return nil
}

//
//func trimHexPrefix(hexString string) (string, error) {
//	if len(hexString) < 2 {
//		return "", fmt.Errorf("invalid hex string %s", hexString)
//	}
//	if hexString[:2] == "0x" {
//		return hexString[2:], nil
//	}
//	return hexString, nil
//}
//
//func encode(b []byte) string {
//	enc := make([]byte, len(b)*2+2)
//	copy(enc, "0x")
//	hex.Encode(enc[2:], b)
//	return string(enc)
//}
//
//// https://github.com/immutable/imx-core-sdk-golang/issues/81
//func getPublicKey(privateKeyInHex string) string {
//	privateKey, _ := trimHexPrefix(privateKeyInHex)
//	privateKeyInEcdsa, _ := crypto.HexToECDSA(privateKey)
//	pubKey := crypto.FromECDSAPub(&privateKeyInEcdsa.PublicKey)
//	return encode(pubKey)
//}
//
//func updateCollection(c *imx.Client, l1signer imx.L1Signer, info *CollectionInformation) {
//	ctx := context.TODO()
//
//	request := api.NewUpdateCollectionRequest()
//	request.Name = &info.CollectionName
//	request.MetadataApiUrl = &info.MetadataUrl
//	request.Description = &info.CollectionName
//
//	response, err := c.UpdateCollection(ctx, l1signer, info.ContractAddress, request)
//	if err != nil {
//		log.Panicf("error in CreateCollection: %v\n", err)
//	}
//
//	val, err := json.MarshalIndent(response, "", "    ")
//	if err != nil {
//		log.Panicf("error in json marshaling: %v\n", err)
//	}
//	log.Println("Created new collection, response: ", string(val))
//
//	// Get the collection details we just created.
//	collectionReponse, err := c.GetCollection(ctx, info.ContractAddress)
//	if err != nil {
//		log.Panicf("error when calling `GetCollection: %v", err)
//	}
//	log.Println("Created Collection Name: ", collectionReponse.Name)
//}
//

//

//
//func getBoolPointer(val bool) *bool {
//	return &val
//}
//
//func getToken(c *imx.Client, address string) {
//	ctx := context.TODO()
//	asset, err := c.GetAsset(ctx, address, "3", getBoolPointer(true))
//	if err != nil {
//		log.Panicf("error in AddMetadataSchemaToCollection: %v\n", err)
//	}
//	val, err := json.MarshalIndent(asset, "", "    ")
//	if err != nil {
//		log.Panicf("error in json marshaling: %v\n", err)
//	}
//	log.Println("Created new metadata, response: ", string(val))
//}
//
//func getMetadata(c *imx.Client) {
//	ctx := context.TODO()
//	meta, err := c.GetMetadataSchema(ctx, "0x4958d0B91412eE2b8D715bF9279DCDB68e33d195")
//	if err != nil {
//		log.Panicf("error in AddMetadataSchemaToCollection: %v\n", err)
//	}
//	val, err := json.MarshalIndent(meta, "", "    ")
//	if err != nil {
//		log.Panicf("error in json marshaling: %v\n", err)
//	}
//	log.Println("Created new metadata, response: ", string(val))
//}
//
//func metadataRefresh(c *imx.Client, l1signer imx.L1Signer) {
//	ctx := context.TODO()
//
//	request := api.NewCreateMetadataRefreshRequest("0x4958d0B91412eE2b8D715bF9279DCDB68e33d195", []string{"1", "2", "3"})
//
//	response, err := c.CreateMetadataRefresh(ctx, l1signer, request)
//	if err != nil {
//		log.Panicf("error in AddMetadataSchemaToCollection: %v\n", err)
//	}
//	val, err := json.MarshalIndent(response, "", "    ")
//	if err != nil {
//		log.Panicf("error in json marshaling: %v\n", err)
//	}
//	log.Println("Created new metadata, response: ", string(val))
//}
//

//

//func getOrders(c *imx.Client) {
//	request := api.ApiListOrdersRequest{}
//	request = request.User("0x1E09BCED9684d94fDCa0b3c7f42F3F21D0d32b4d")
//	response, err := c.ListOrders(&request)
//	if err != nil {
//		log.Panicf("error in CreateOrder: %v", err)
//	}
//
//	createOrderResponseStr, err := prettyStruct(response)
//	if err != nil {
//		log.Panic(err)
//	}
//	log.Printf("CreateOrder response:\n%v\n", createOrderResponseStr)
//}
//
//func createEthDeposit(c *imx.Client, l1signer imx.L1Signer, info *CreateDeposit) {
//	ctx := context.TODO()
//	// Eth Deposit
//	ethAmountInWei, err := strconv.ParseUint(info.DepositAmountWei, 10, 64)
//	if err != nil {
//		log.Panicf("error in converting ethAmountInWei from string to int: %v\n", err)
//	}
//
//	transaction, err := imx.NewETHDeposit(ethAmountInWei).Deposit(ctx, c, l1signer, nil)
//	if err != nil {
//		log.Panicf("Eth deposit: %v", err)
//	}
//	log.Println("Eth Deposit transaction hash:", transaction.Hash())
//}
//
//func createTrade(c *imx.Client, l1signer imx.L1Signer, l2signer imx.L2Signer, orderID int32) {
//	ctx := context.TODO()
//	tradeRequest := api.GetSignableTradeRequest{
//		Fees:    nil,
//		OrderId: orderID,
//	}
//	tradeRequest.SetExpirationTimestamp(0)
//	tradeResponse, err := c.CreateTrade(ctx, l1signer, l2signer, tradeRequest)
//
//	if err != nil {
//		log.Panicf("error calling trades workflow: %v", err)
//	}
//
//	val, _ := json.MarshalIndent(tradeResponse, "", "  ")
//	log.Printf("trade response:\n%s\n", val)
//}
//
//func createEthWithdrawal(c *imx.Client, l1signer imx.L1Signer, l2signer imx.L2Signer, amount string) {
//	ctx := context.TODO()
//	ethAmountInWei, err := strconv.ParseUint(amount, 10, 64)
//	if err != nil {
//		log.Panicf("error in converting ethAmountInWei from string to int: %v\n", err)
//	}
//
//	withdrawalRequest := api.GetSignableWithdrawalRequest{
//		Amount: strconv.FormatUint(ethAmountInWei, 10),
//		Token:  imx.SignableETHToken(),
//	}
//
//	response, err := c.PrepareWithdrawal(ctx, l1signer, l2signer, withdrawalRequest)
//	if err != nil {
//		log.Panicf("error calling PrepareWithdrawal workflow: %v", err)
//	}
//	val, _ := json.MarshalIndent(response, "", "  ")
//	log.Printf("response:\n%s\n", val)
//
//}
//
//func completeEthWithdrawal(c *imx.Client, l1signer imx.L1Signer, l2signer imx.L2Signer, withdrawalId int32) {
//	ctx := context.TODO()
//	for {
//		getWithdrawalResponse, err := c.GetWithdrawal(ctx, strconv.FormatInt(int64(withdrawalId), 10))
//		if err != nil {
//			log.Panicf("error calling GetWithdrawal: %v", err)
//		}
//		val, _ := json.MarshalIndent(getWithdrawalResponse, "", "  ")
//		log.Printf("response:\n%s\n", val)
//
//		if getWithdrawalResponse.RollupStatus == "confirmed" {
//			break
//		}
//		time.Sleep(5 * time.Minute)
//	}
//
//	ethWithdrawal := imx.NewEthWithdrawal()
//	transaction, err := ethWithdrawal.CompleteWithdrawal(ctx, c, l1signer, l2signer.GetPublicKey(), nil)
//	if err != nil {
//		log.Panicf("error calling withdrawalsWorkflow.CompleteEthWithdrawal workflow: %v", err)
//	}
//	log.Println("transaction hash:", transaction.Hash())
//}
//
//func listCollections(c *imx.Client, keyword string) {
//	request := api.ApiListCollectionsRequest{}
//	request = request.Keyword(keyword)
//	request = request.PageSize(10)
//	response, err := c.ListCollections(&request)
//	if err != nil {
//		log.Panicf("error calling PrepareWithdrawal workflow: %v", err)
//	}
//	val, _ := json.MarshalIndent(response, "", "  ")
//	log.Printf("response:\n%s\n", val)
//}
