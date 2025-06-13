package sms

const _apiKeyHeader = "X-Api-Key"

func (s *sms) headers() map[string]string {
	return map[string]string{
		_apiKeyHeader: s.config.GetString("sms.token"),
	}
}
