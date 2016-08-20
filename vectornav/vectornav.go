package vectornav

import (
	"fmt"
	"strconv"
	"strings"
)

// YMRData represent $VNYMR sentence from Vectornav IMU
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
type YMRData struct {
	Yaw   float64
	Pitch float64
	Roll  float64

	MagX float64
	MagY float64
	MagZ float64

	AccelX float64
	AccelY float64
	AccelZ float64

	GyroX float64
	GyroY float64
	GyroZ float64
}

func parseYMR(s string) (*YMRData, error) {
	fieldsChsum := strings.Split(s, "*")
	// separate the checksum
	if len(fieldsChsum) != 2 {
		return nil, fmt.Errorf("Unable to split checksum, got %d parts", len(fieldsChsum))
	}

	// calculate the checksum (XOR of bytes between $ and * - without these)
	var chsum int64
	toChsStr := fieldsChsum[0][1:len(fieldsChsum[0])]
	for _, c := range []byte(toChsStr) {
		chsum ^= int64(c)
	}
	msgChsum, err := strconv.ParseInt(fieldsChsum[1][0:2], 16, 64)
	if err != nil {
		return nil, fmt.Errorf("Checksum parse error: %v", err)
	}
	if msgChsum != chsum {
		return nil, fmt.Errorf("Checksum error, got: %x, expected %x", msgChsum, chsum)
	}

	// split fields
	fields := strings.Split(fieldsChsum[0], ",")
	if fields[0] != "$VNYMR" {
		return nil, fmt.Errorf("No VNYMR message, found %s", fields[0])
	}
	if len(fields) != 13 {
		return nil, fmt.Errorf("Unable to split fields, got %d parts", len(fields))
	}

	var data YMRData
	// and now fill the fields
	data.Yaw, err = strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, fmt.Errorf("Yaw parse error: %v", err)
	}
	data.Pitch, err = strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, fmt.Errorf("Pitch parse error: %v", err)
	}
	data.Roll, err = strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return nil, fmt.Errorf("Roll parse error: %v", err)
	}
	data.MagX, err = strconv.ParseFloat(fields[4], 64)
	if err != nil {
		return nil, fmt.Errorf("MagX parse error: %v", err)
	}
	data.MagY, err = strconv.ParseFloat(fields[5], 64)
	if err != nil {
		return nil, fmt.Errorf("MagY parse error: %v", err)
	}
	data.MagZ, err = strconv.ParseFloat(fields[6], 64)
	if err != nil {
		return nil, fmt.Errorf("MagZ parse error: %v", err)
	}
	data.AccelX, err = strconv.ParseFloat(fields[7], 64)
	if err != nil {
		return nil, fmt.Errorf("AccelX parse error: %v", err)
	}
	data.AccelY, err = strconv.ParseFloat(fields[8], 64)
	if err != nil {
		return nil, fmt.Errorf("AccelY parse error: %v", err)
	}
	data.AccelZ, err = strconv.ParseFloat(fields[9], 64)
	if err != nil {
		return nil, fmt.Errorf("AccelZ parse error: %v", err)
	}
	data.GyroX, err = strconv.ParseFloat(fields[10], 64)
	if err != nil {
		return nil, fmt.Errorf("GyroX parse error: %v", err)
	}
	data.GyroY, err = strconv.ParseFloat(fields[11], 64)
	if err != nil {
		return nil, fmt.Errorf("GyroY parse error: %v", err)
	}
	data.GyroZ, err = strconv.ParseFloat(fields[12], 64)
	if err != nil {
		return nil, fmt.Errorf("GyroZ parse error: %v", err)
	}

	return &data, nil
}
