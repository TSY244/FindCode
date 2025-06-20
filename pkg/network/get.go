package network

//func GetRequest(target string, header Header) (string, error) {
//
//	req, _ := http.NewRequest(Get, target, nil)
//	for key, value := range header {
//		req.Header.Set(key, value)
//	}
//
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		return "", fmt.Errorf("http.RequestWithJson -> %v", err)
//	}
//	defer resp.Body.Close()
//
//	buf := new(bytes.Buffer)
//	buf.ReadFrom(resp.Body)
//	return buf.String(), nil
//}
