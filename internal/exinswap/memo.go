package exinswap

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// TradeMemoV2 builds ExinSwap V2 trade memo.
// Format before base64: ACTION$target_asset_uuid$min_out$latest_exec_time$route
// Optional fields can be omitted or set to "0".
func TradeMemoV2(targetAssetUUID string, minOut string, latestExec *time.Time, route string) (string, error) {
	if targetAssetUUID == "" {
		return "", fmt.Errorf("missing targetAssetUUID")
	}
	fields := []string{"0", targetAssetUUID}
	if minOut != "" {
		fields = append(fields, minOut)
	}
	if latestExec != nil {
		fields = append(fields, fmt.Sprintf("%d", latestExec.UTC().Unix()))
	}
	if route != "" {
		// ensure latestExec present if route specified? doc allows optional; keep as-is.
		fields = append(fields, route)
	}
	raw := strings.Join(fields, "$")
	return base64.StdEncoding.EncodeToString([]byte(raw)), nil
}
