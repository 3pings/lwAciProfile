package main

import (
	"fmt"
	"github.com/3pings/acigo/aci"
	"log"
	"os"
	//"os"
)

func main() {

	//Get Hostname
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	//Print hostname for testing
	fmt.Println("hostname:", host)

	//aciClient Info
	//host := "pod01"
	mysqlContract := host + "-mysql"

	allAllContract := host + "-l3out-allow-all"

	a, errLogin := login(false)
	if errLogin != nil {
		log.Printf("exiting: %v", errLogin)
		return
	}

	defer logout(a)

	//Add kube-nodes health-check provided contract
	errAddKubeNodesHC := a.EPGContractProvidedAdd(host, "kubernetes", "kube-nodes", "health-check")
	if errAddKubeNodesHC != nil {
		fmt.Printf("Kube-Nodes provided health-check contract add error: %v\n", errAddKubeNodesHC)
		return
	}

	//Add kube-default hx-connect consume contract
	errAddKubeDefaultHX := a.EPGContractConsumedAdd(host, "kubernetes", "kube-default", "hx-connect")
	if errAddKubeDefaultHX != nil {
		fmt.Printf("Kube-Default hx-connect consume contract add error: %v\n", errAddKubeDefaultHX)
		return
	}

	//Add kube-system hx-connect consume contract
	errAddKubeSystemDNS := a.EPGContractConsumedAdd(host, "kubernetes", "kube-system", "dns")
	if errAddKubeSystemDNS != nil {
		fmt.Printf("frontend Consume API contract add error: %v\n", errAddKubeSystemDNS)
		return
	}

	//Add Local API Contract
	errAddLocalAPIContract := a.ContractAdd(host, "api", "tenant", "")
	if errAddLocalAPIContract != nil {
		fmt.Printf("Error creating API contract add error: %v\n", errAddLocalAPIContract)
		return
	}

	//ADD API Contract Subject
	errAPIContractSubjectAdd := a.ContractSubjectAdd(host, "api", "api", "true", true, "")
	if errAPIContractSubjectAdd != nil {
		fmt.Printf("Error adding Contract Subject: %v\n", errAPIContractSubjectAdd)
		return
	}

	//Create Filter
	errAddAPIFilter := a.FilterAdd(host, "api-filter", "")
	if errAddAPIFilter != nil {
		fmt.Printf("Error adding filter to contract: %v\n", errAddAPIFilter)
		return
	}

	//Add Entry into Filter for API
	errAPIFilterEntryAdd := a.FilterEntryAdd(host, "api-filter", "api", "ip", "tcp", "unspecified", "unspecified", "5000", "5000")
	if errAPIFilterEntryAdd != nil {
		fmt.Printf("Error adding filter entry for API: %v\n", errAPIFilterEntryAdd)
		return
	}
	//Add Filter to Subject for API
	errAddAPISubjectApplyBoth := a.SubjectFilterBothAdd(host, "api", "api", "api-filter")
	if errAddAPISubjectApplyBoth != nil {
		fmt.Printf("Error adding filter entry for API: %v\n", errAddAPISubjectApplyBoth)
		return
	}

	//Add application profile for livewall
	errAddAppProfile := a.ApplicationProfileAdd(host, "livewall", "")
	if errAddAppProfile != nil {
		fmt.Printf("app profile add error: %v\n", errAddAppProfile)
		return
	}

	//Logic for FrontEnd EPG
	//Add frontend EPG to livewall app profile
	errAddFrontEnd := a.ApplicationEPGAdd(host, "livewall", "kube-pod-bd", "frontend", "")
	if errAddFrontEnd != nil {
		fmt.Printf("frontend epg add error: %v\n", errAddFrontEnd)
		return
	}

	//Add VMM Domain to frontend EPG
	errFrontEndEPGVMMAdd := a.EPGVMMAdd("Kubernetes", host, "livewall", "frontend")
	if errFrontEndEPGVMMAdd != nil {
		fmt.Printf("Frontend VMM Association Failed: %v\n", errFrontEndEPGVMMAdd)
		return
	}

	errAddFrontendInheritContract := a.EPGContractInheritanceAdd(host, "livewall", "frontend")
	if errAddFrontendInheritContract != nil {
		fmt.Printf("FrontEnd consume API contract add error: %v\n", errAddFrontendInheritContract)
		return
	}
	//Add FrontEnd API consume contract
	errAddFrontEndAPI := a.EPGContractConsumedAdd(host, "livewall", "frontend", "api")
	if errAddFrontEndAPI != nil {
		fmt.Printf("FrontEnd consume API contract add error: %v\n", errAddFrontEndAPI)
		return
	}

	//Logic for API EPG
	//Add API EPG to livewall app profile
	errAddAPI := a.ApplicationEPGAdd(host, "livewall", "kube-pod-bd", "api", "")
	if errAddAPI != nil {
		fmt.Printf("API epg add error: %v\n", errAddAPI)
		return
	}

	//Add VMM Domain to API EPG
	errAPIEPGVMMAdd := a.EPGVMMAdd("Kubernetes", host, "livewall", "api")
	if errAPIEPGVMMAdd != nil {
		fmt.Printf("API VMM Association Failed: %v\n", errAPIEPGVMMAdd)
		return
	}

	//Add API Inherited Contract
	errAddApiInheritContract := a.EPGContractInheritanceAdd(host, "livewall", "api")
	if errAddApiInheritContract != nil {
		fmt.Printf("API consume API contract add error: %v\n", errAddApiInheritContract)
		return
	}
	//Add API API provided contract
	errAddAPIAPI := a.EPGContractProvidedAdd(host, "livewall", "api", "api")
	if errAddAPIAPI != nil {
		fmt.Printf("API provided health-check contract add error: %v\n", errAddAPIAPI)
		return
	}

	//Add API MySql consume contract
	errAddAPIMySql := a.EPGContractConsumedAdd(host, "livewall", "api", mysqlContract)
	if errAddAPIMySql != nil {
		fmt.Printf("API Consume MySql contract add error: %v\n", errAddAPIMySql)
		return
	}

	//Logic for collector EPG
	//Add collector EPG to livewall app profile
	errAddCollector := a.ApplicationEPGAdd(host, "livewall", "kube-pod-bd", "collector", "")
	if errAddCollector != nil {
		fmt.Printf("Collector epg add error: %v\n", errAddCollector)
		return
	}

	//Add VMM Domain to collector EPG
	errCollectorEPGVMMAdd := a.EPGVMMAdd("Kubernetes", host, "livewall", "collector")
	if errCollectorEPGVMMAdd != nil {
		fmt.Printf("Collector VMM Association Failed: %v\n", errCollectorEPGVMMAdd)
		return
	}

	//Add Collector Inherited Contract
	errAddCollectorInheritContract := a.EPGContractInheritanceAdd(host, "livewall", "collector")
	if errAddCollectorInheritContract != nil {
		fmt.Printf("Collector consume API contract add error: %v\n", errAddCollectorInheritContract)
		return
	}

	//Add collector MySql consume contract
	errAddCollectorMySql := a.EPGContractConsumedAdd(host, "livewall", "collector", mysqlContract)
	if errAddCollectorMySql != nil {
		fmt.Printf("Collector Consume MySql contract add error: %v\n", errAddCollectorMySql)
		return
	}

	//Add collector l3out consume contract
	errAddCollectorL3out := a.EPGContractConsumedAdd(host, "livewall", "collector", allAllContract)
	if errAddCollectorL3out != nil {
		fmt.Printf("Collector Consume L3Out contract add error: %v\n", errAddCollectorL3out)
		return
	}
	fmt.Println("Livewall Profile Successfully added to " + host + "tenant.")
}

func login(debug bool) (*aci.Client, error) {
	//Apic login information variables
	apicHosts := []string{""}
	apicUser := ""
	apicPass := ""
	a, errNew := aci.New(aci.ClientOptions{Hosts: apicHosts, User: apicUser, Pass: apicPass, Debug: debug})
	if errNew != nil {
		return nil, fmt.Errorf("login new client error: %v", errNew)
	}

	errLogin := a.Login()
	if errLogin != nil {
		return nil, fmt.Errorf("login error: %v", errLogin)
	}

	return a, nil
}

func logout(a *aci.Client) {
	errLogout := a.Logout()
	if errLogout != nil {
		log.Printf("logout error: %v", errLogout)
		return
	}
}
