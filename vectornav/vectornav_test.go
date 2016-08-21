package vectornav

import "testing"

func TestParsingCorrect(t *testing.T) {
	var data YMRDataFull
	err := ParseYMR("$VNYMR,+104.977,+004.548,-001.276,-00.8012,-02.7376,+01.0070,+00.837,+00.235,-10.414,-00.002081,-00.001151,+00.002113*61\r\n", &data)
	if err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}
	if data.Yaw != 104.977 {
		t.Errorf("Yaw not matched: got %f, expected %f", data.Yaw, 104.977)
	}
	if data.Pitch != 4.548 {
		t.Errorf("Pitch not matched: got %f, expected %f", data.Pitch, 4.548)
	}
	if data.Roll != -1.276 {
		t.Errorf("Roll not matched: got %f, expected %f", data.Roll, -1.276)
	}
	if data.MagX != -0.8012 {
		t.Errorf("MagX not matched: got %f, expected %f", data.MagX, -0.8012)
	}
	if data.MagY != -2.7376 {
		t.Errorf("MagY not matched: got %f, expected %f", data.MagY, -2.7376)
	}
	if data.MagZ != 1.007 {
		t.Errorf("MagZ not matched: got %f, expected %f", data.MagZ, 1.007)
	}
	if data.AccelX != 0.837 {
		t.Errorf("AccelX not matched: got %f, expected %f", data.AccelX, 0.837)
	}
	if data.AccelY != 0.235 {
		t.Errorf("AccelY not matched: got %f, expected %f", data.AccelY, 0.235)
	}
	if data.AccelZ != -10.414 {
		t.Errorf("AccelZ not matched: got %f, expected %f", data.AccelZ, -10.414)
	}
	if data.GyroX != -0.002081 {
		t.Errorf("GyroX not matched: got %f, expected %f", data.GyroX, -0.002081)
	}
	if data.GyroY != -0.001151 {
		t.Errorf("GyroY not matched: got %f, expected %f", data.GyroY, -0.001151)
	}
	if data.GyroZ != 0.002113 {
		t.Errorf("GyroZ not matched: got %f, expected %f", data.GyroZ, 0.002113)
	}
}

func TestParsingInvalidChsum(t *testing.T) {
	var data YMRDataFull
	err := ParseYMR("$VNYMR,+104.977,+004.548,-001.276,-00.8012,-02.7376,+01.0070,+00.837,+00.235,-10.414,-00.002081,-00.001151,+00.002113*60\r\n", &data)
	if err == nil {
		t.Errorf("Parsing did not failed")
	}

	if err.Error() != "Checksum error, got: 60, expected 61" {
		t.Errorf("Unexpected error message, got: %s", err.Error())
	}
}
