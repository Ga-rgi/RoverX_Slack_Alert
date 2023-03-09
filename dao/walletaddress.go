package dao

type WhitelistedAddresses struct {
	Address          string
	IsWhitelisted    bool
	PartnerCommunity string
}

var whitelistAddresses []WhitelistedAddresses

func CheckAddressWhitelisted(address string) bool {
	WhitelistAddresses()
	for _, item := range whitelistAddresses {
		if item.Address == address {
			return item.IsWhitelisted
		}
	}
	return false
}

// user addresses
func WhitelistAddresses() {
	//whitelisted
	whitelistAddresses = append(whitelistAddresses, WhitelistedAddresses{Address: "0x4A906262CFE6B4de05A3E0b890Bf8eb4a4c2f30A", IsWhitelisted: true, PartnerCommunity: "ZenAcademy"})

	//NO COMMUNITY
	whitelistAddresses = append(whitelistAddresses, WhitelistedAddresses{Address: "0x4A906262CFE6B4de05A3E0b890Bf8eb4a4c2f30B", IsWhitelisted: true, PartnerCommunity: ""})
}

func GetPartnerCommunity(address string) string {
	for _, item := range whitelistAddresses {
		if item.Address == address {
			return item.PartnerCommunity
		}
	}
	return ""
}
