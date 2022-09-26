package g

import "lihood/conf"

func ChainCallback(oid string) string {
	return conf.Instance.Server.Domain + "/api/v1/product/chain/callback/" + oid
}
