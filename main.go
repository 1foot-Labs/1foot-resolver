package main

import (
	"encoding/json"
	"net/http"
	"sort"
    "time"

	"fmt"

	"os/exec"

    "math/big"
	"resolver/models"
    "resolver/HTLC"
    "resolver/relayer_communication"
)

func lookForActiveOrders() {
	resp, err := http.Get("http://localhost:3002/api/active-orders")
	if err != nil {
		fmt.Println("❌ Error calling API:", err)
		return
	}
	defer resp.Body.Close()

	var orders []models.ActiveOrder
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		fmt.Println("❌ Error decoding response:", err)
		return
	}

	if len(orders) == 0 {
		fmt.Println("⚠️ No active orders found.")
		return
	}

	// Sort by createdAt (latest last)
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].CreatedAt < orders[j].CreatedAt
	})

	latest := orders[len(orders)-1]

	fmt.Println("✅ Latest Order:")
	fmt.Printf("ID: %s | ETH Addr: %s | sha256: %s\n", latest.ID, latest.MakerAddress, latest.Sha256)

	_, p2shAddress := HTLC.CreateHTLCForBitcoin(latest.Sha256, latest.PubKey)

    clientURL := "http://localhost:8545" // your Ethereum node
	privateKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"  // resolver pvt key
	contractAddressHex := "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
	receiverHex := latest.MakerAddress 
	secretHex := latest.Sha256       
	timelockSeconds := int64(3600 * 24)  
	amountInWei := big.NewInt(int64(latest.AmountToReceive))     

	// Call the function
	txHash, err := HTLC.CreateHTLCForEthereum(
		clientURL,
		privateKeyHex,
		contractAddressHex,
		receiverHex,
		secretHex,
		timelockSeconds,
		amountInWei,
	)
	if err != nil {
		fmt.Println("Error creating HTLC:", err)
		return
	}

	fmt.Println("HTLC created with tx hash:", txHash)

    request := relayer_communication.FulfillOrderRequest{
		OrderID:        latest.ID,
		TakerAddress:   "",
		EthHTLCAddress: contractAddressHex,
		BtcHTLCAddress: p2shAddress,
	}

	resp, err = relayer_communication.FulfillOrder(request)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.StatusCode)
}

func getCurrentBlockHeight() (int64, error) {
	cmd := exec.Command("bitcoin-cli", "-regtest", "getblockcount")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	var height int64
	_, err = fmt.Sscanf(string(output), "%d", &height)
	if err != nil {
		return 0, err
	}
	return height, nil
}

func main() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	lookForActiveOrders()

	for range ticker.C {
		lookForActiveOrders()
	}
	// HTLC.CreateHTLCForBitcoin("6b6a5987b7a4cbbf2310cb1e785df165f83a6e44",
	// 	"022514f3c0d22eac4d45ecc6ed9fb17fa44cebb88d590b79ca834b20a552f9bb67")

    // clientURL := "http://127.0.0.1:8545" // your Ethereum node
	// privateKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"       // your private key (without 0x prefix or with)
	// contractAddressHex := "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"
	// receiverHex := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266" // receiver's Ethereum address
	// secretSHA256Hex := "fc56dbc6d4652b315b86b71c8d688c1ccdea9c5f1fd07763d2659fde2e2fc49a"       // hex-encoded secret (32 bytes)
	// timelockSeconds := int64(3600 * 24)       // 24 hour
	// amountInWei := big.NewInt(1e18)      // 1 ETH in wei

	// // Call the function
	// txHash, err := HTLC.CreateHTLCForEthereum(
	// 	clientURL,
	// 	privateKeyHex,
	// 	contractAddressHex,
	// 	receiverHex,
	// 	secretSHA256Hex,
	// 	timelockSeconds,
	// 	amountInWei,
	// )
	// if err != nil {
	// 	fmt.Println("Error creating HTLC:", err)
	// 	return
	// }

	// fmt.Println("HTLC created with tx hash:", txHash)
}
