package vectornav

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// YMRData represent basic data of $VNYMR sentence from Vectornav IMU
type YMRData struct {
	Yaw   float64 `json:"yaw"`
	Pitch float64 `json:"pitch"`
	Roll  float64 `json:"roll"`
}

// YMRDataFull represent $VNYMR sentence from Vectornav IMU
// Sample: "$VNYMR,+104.977,+004.548,-001.276,-00.8012,-02.7376,+01.0070,+00.837,+00.235,-10.414,-00.002081,-00.001151,+00.002113*61\r\n"
// Field # Meaning Units Printf/Scanf format
// 1 Yaw Angle degrees %+07.2f
// 2 Pitch Angle degrees %+07.2f
// 3 Roll Angle degrees %+07.2f
// 4 X-Axis Magnetic N/A %+07.4f
// 5 Y-Axis Magnetic N/A %+07.4f
// 6 Z-Axis Magnetic N/A %+07.4f
// 7 X-Axis Acceleration m/s^2 %+07.3f
// 8 Y-Axis Acceleration m/s^2 %+07.3f
// 9 Z-Axis Acceleration m/s^2 %+07.3f
// 10 X-Axis Angular Rate rad/s %+08.4f
// 11 Y-Axis Angular Rate rad/s %+08.4f
// 12 Z-Axis Angular Rate rad/s %+08.4f
type YMRDataFull struct {
	YMRData

	MagX float64 `json:"magX"`
	MagY float64 `json:"magY"`
	MagZ float64 `json:"magZ"`

	AccelX float64 `json:"accelX"`
	AccelY float64 `json:"accelY"`
	AccelZ float64 `json:"accelZ"`

	GyroX float64 `json:"gyroX"`
	GyroY float64 `json:"gyroY"`
	GyroZ float64 `json:"gyroZ"`
}

// ParseYMR reads the string representation of VNYMR sentence and fills data structure
func ParseYMR(s string, data *YMRDataFull) error {
	fieldsChsum := strings.Split(s, "*")
	// separate the checksum
	if len(fieldsChsum) != 2 {
		return fmt.Errorf("Unable to split checksum, got %d parts", len(fieldsChsum))
	}

	// calculate the checksum (XOR of bytes between $ and * - without these two)
	var chsum int64
	toChsStr := fieldsChsum[0][1:len(fieldsChsum[0])]
	for _, c := range []byte(toChsStr) {
		chsum ^= int64(c)
	}
	msgChsum, err := strconv.ParseInt(fieldsChsum[1][0:2], 16, 64) // trim the newline
	if err != nil {
		return fmt.Errorf("Checksum parse error: %v", err)
	}
	if msgChsum != chsum {
		return fmt.Errorf("Checksum error, got: %x, expected %x", msgChsum, chsum)
	}

	// split fields
	fields := strings.Split(fieldsChsum[0], ",")
	if fields[0] != "$VNYMR" {
		return fmt.Errorf("No VNYMR message, found %s", fields[0])
	}
	if len(fields) != 13 {
		return fmt.Errorf("Unable to split fields, got %d parts", len(fields))
	}

	// and now fill the fields
	data.Yaw, err = strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return fmt.Errorf("Yaw parse error: %v", err)
	}
	data.Pitch, err = strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return fmt.Errorf("Pitch parse error: %v", err)
	}
	data.Roll, err = strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return fmt.Errorf("Roll parse error: %v", err)
	}
	data.MagX, err = strconv.ParseFloat(fields[4], 64)
	if err != nil {
		return fmt.Errorf("MagX parse error: %v", err)
	}
	data.MagY, err = strconv.ParseFloat(fields[5], 64)
	if err != nil {
		return fmt.Errorf("MagY parse error: %v", err)
	}
	data.MagZ, err = strconv.ParseFloat(fields[6], 64)
	if err != nil {
		return fmt.Errorf("MagZ parse error: %v", err)
	}
	data.AccelX, err = strconv.ParseFloat(fields[7], 64)
	if err != nil {
		return fmt.Errorf("AccelX parse error: %v", err)
	}
	data.AccelY, err = strconv.ParseFloat(fields[8], 64)
	if err != nil {
		return fmt.Errorf("AccelY parse error: %v", err)
	}
	data.AccelZ, err = strconv.ParseFloat(fields[9], 64)
	if err != nil {
		return fmt.Errorf("AccelZ parse error: %v", err)
	}
	data.GyroX, err = strconv.ParseFloat(fields[10], 64)
	if err != nil {
		return fmt.Errorf("GyroX parse error: %v", err)
	}
	data.GyroY, err = strconv.ParseFloat(fields[11], 64)
	if err != nil {
		return fmt.Errorf("GyroY parse error: %v", err)
	}
	data.GyroZ, err = strconv.ParseFloat(fields[12], 64)
	if err != nil {
		return fmt.Errorf("GyroZ parse error: %v", err)
	}

	return nil
}

// GetVectornavHTTPHandler returns handler returning IMU data in JSON format
func GetVectornavHTTPHandler(data *YMRData) func(w http.ResponseWriter, req *http.Request) {
	fn := func(w http.ResponseWriter, req *http.Request) {
		js, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
	return fn
}
