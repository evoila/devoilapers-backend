package provider

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/postgres"
	"OperatorAutomation/pkg/postgres/dtos/action_dtos"
	"OperatorAutomation/pkg/postgres/dtos/provider_dtos"
	"OperatorAutomation/test/integration_tests/common_test"
	unit_test "OperatorAutomation/test/unit_tests/common_test"
	"encoding/base64"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

const testCaBase64 = "Q2VydGlmaWNhdGU6CiAgICBEYXRhOgogICAgICAgIFZlcnNpb246IDEgKDB4MCkKICAgICAgICBTZXJpYWwgTnVtYmVyOgogICAgICAgICAgICA0NzpiOTpiYTo3ZTo4YjpjZTpiNzozNToyYjo1ODpmODpiYzo4Mjo1Yjo0ODo1ZjowYzo4NjplMjozMAogICAgICAgIFNpZ25hdHVyZSBBbGdvcml0aG06IE5VTEwKICAgICAgICBJc3N1ZXI6IENOID0gc29tZS5jZXJ0LmRlCiAgICAgICAgVmFsaWRpdHkKICAgICAgICAgICAgTm90IEJlZm9yZTogTWFyIDIzIDEzOjM0OjE2IDIwMjEgR01UCiAgICAgICAgICAgIE5vdCBBZnRlciA6IE1hciAyMSAxMzozNDoxNiAyMDMxIEdNVAogICAgICAgIFN1YmplY3Q6IENOID0gc29tZS5jZXJ0LmRlCiAgICAgICAgU3ViamVjdCBQdWJsaWMgS2V5IEluZm86CiAgICAgICAgICAgIFB1YmxpYyBLZXkgQWxnb3JpdGhtOiByc2FFbmNyeXB0aW9uCiAgICAgICAgICAgICAgICBSU0EgUHVibGljLUtleTogKDIwNDggYml0KQogICAgICAgICAgICAgICAgTW9kdWx1czoKICAgICAgICAgICAgICAgICAgICAwMDpjNDoyYjozYTplNjo3NjpmNTo5OTpkZTpjODo2ZjowZDpkODoyZDphNDoKICAgICAgICAgICAgICAgICAgICBmYToyYzpmNDpkYTplNjphNzoyMTo2ODpjNzpjZTowNzpmMDo4MjpjNzoyODoKICAgICAgICAgICAgICAgICAgICA3NDo4MzpmYTpkNToxNTphYzpjMDo5NTo5YToyMzpjMjo0NTo2YjpmODo5MToKICAgICAgICAgICAgICAgICAgICAzMjoxNDowNTo2NToxOTo3ODozZjo1Nzo1NDo5NTo0NDpiMTowNTo0ODoyYjoKICAgICAgICAgICAgICAgICAgICBjMjphNDplYTo1ODowMTo2OToxNjo4OTo4OTo0ZDplOTo4MzoxZjpmOToyYToKICAgICAgICAgICAgICAgICAgICA3MjphZTpjNjpiYTpkNjo3ZTo5NDo3YTozYTo2NDo1ZjoxZjpiYTozODoyNzoKICAgICAgICAgICAgICAgICAgICBhMjo0NTpiMjowNDpjZTo5YTo1MzphZTo2NzpkODplMjo2MDozMDpiNDo3MjoKICAgICAgICAgICAgICAgICAgICBiNjoxYjplZDpmMjpkZTo0NDpiYzo0Njo1Mzo2ZDo3MzplMTo4OTpkODozZjoKICAgICAgICAgICAgICAgICAgICA1ZTowNTowYjplZTpkYTozMTowNzpmZDo0MDplNDo5Yzo3YTo2Mjo5ZTo3YToKICAgICAgICAgICAgICAgICAgICBmNTo4OTowNTo4MTowZjo1ZjoyNTo0MTplMDo0MzpjNDo3Njo4NDoxMDpjNDoKICAgICAgICAgICAgICAgICAgICA2ZjplODo4NDpiNTphNjpmZTpmZDplNTowMzoyMzplNDo4NDo3Yzo2ZTo4MToKICAgICAgICAgICAgICAgICAgICA3Yzo4MTo0ODo2ODo0OTo1NDoyZjozYTpiMzowNTpkNzpkNjphMjpjYTpiZToKICAgICAgICAgICAgICAgICAgICBlYjphYzo2ZjpjNDowOTowNTpjOTo5ZjpjNTo4NDphZjo5Mzo3YToyZjphMToKICAgICAgICAgICAgICAgICAgICAxNzoxNDowMjoxNDpkMjowNjpiNjpkMTphNjpkMTowODpjMDo4OTphNDphNjoKICAgICAgICAgICAgICAgICAgICAwYzpiNTo2MjpjNDo2MTpjYjo3ODpjYzpiNjoxYToyODoxZjoyMjpjNDo3MzoKICAgICAgICAgICAgICAgICAgICBhMTplYTphZDpkOTpkYzoyODoyZDo5MTowMTplOTpmMzpkNTpiMToyOTpiNDoKICAgICAgICAgICAgICAgICAgICBlNjo4ODpkNjo0MjoyNzo0NTo1MzoyMzo5MDoxNDoxODoxYjplOTo3Nzo5MzoKICAgICAgICAgICAgICAgICAgICAyYjo5ZAogICAgICAgICAgICAgICAgRXhwb25lbnQ6IDY1NTM3ICgweDEwMDAxKQogICAgU2lnbmF0dXJlIEFsZ29yaXRobTogTlVMTAotLS0tLUJFR0lOIENFUlRJRklDQVRFLS0tLS0KTUlJRER6Q0NBZmVnQXdJQkFnSVVSN202Zm92T3R6VXJXUGk4Z2x0SVh3eUc0akF3RFFZSktvWklodmNOQVFFTApCUUF3RnpFVk1CTUdBMVVFQXd3TWMyOXRaUzVqWlhKMExtUmxNQjRYRFRJeE1ETXlNekV6TXpReE5sb1hEVE14Ck1ETXlNVEV6TXpReE5sb3dGekVWTUJNR0ExVUVBd3dNYzI5dFpTNWpaWEowTG1SbE1JSUJJakFOQmdrcWhraUcKOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQXhDczY1bmIxbWQ3SWJ3M1lMYVQ2TFBUYTVxY2hhTWZPQi9DQwp4eWgwZy9yVkZhekFsWm9qd2tWcitKRXlGQVZsR1hnL1YxU1ZSTEVGU0N2Q3BPcFlBV2tXaVlsTjZZTWYrU3B5CnJzYTYxbjZVZWpwa1h4KzZPQ2VpUmJJRXpwcFRybWZZNG1Bd3RISzJHKzN5M2tTOFJsTnRjK0dKMkQ5ZUJRdnUKMmpFSC9VRGtuSHBpbm5yMWlRV0JEMThsUWVCRHhIYUVFTVJ2NklTMXB2Nzk1UU1qNUlSOGJvRjhnVWhvU1ZRdgpPck1GMTlhaXlyN3JyRy9FQ1FYSm44V0VyNU42TDZFWEZBSVUwZ2EyMGFiUkNNQ0pwS1lNdFdMRVljdDR6TFlhCktCOGl4SE9oNnEzWjNDZ3RrUUhwODlXeEtiVG1pTlpDSjBWVEk1QVVHQnZwZDVNcm5RSURBUUFCbzFNd1VUQWQKQmdOVkhRNEVGZ1FVc0tNK2pXbmVSajJkaUFLdnk1THNhOWR3d2pJd0h3WURWUjBqQkJnd0ZvQVVzS00ralduZQpSajJkaUFLdnk1THNhOWR3d2pJd0R3WURWUjBUQVFIL0JBVXdBd0VCL3pBTkJna3Foa2lHOXcwQkFRc0ZBQU9DCkFRRUFyUTVqK043UjJDQS9kaEJNMVFHWENoa05XTXAybzlKRVJYekw3R3RYNDNDbXRzbnZpMkVISU1iL3ZHemwKc0FDbGI4SkQzb0NES0pabzUzWXlBTUdUV2tuZEZhc3Q5ejJhQ2VkS1ZKL0hrcnRSekVwOEl4Z0JQZzczeTgwZQpPcEcveEIrSUhVdlJMYWVBYWR0UHo0SGdJOU5VTG4wUkFGZk5Yc0FzZmhTUk5LU0svRksxQVNsSTh2Y3ltY1IxClc5cXlOVTI5S2UrWisvb29qaE5HQVdnd1dmZ3NzcE1GN2o5aU9Wd0txNUhrTjRSN3lQNGZuc05HMnFhWExkMFEKeGZLditMNG9SZEpnNkc4aDNseDBoM0tzcXhwYm8zV1lUd0RVOUU0dXdseDZCSmVFc2p6V29adTRiUXVPbjh2ZAorYXZOMGFYOEpyV0NUYmY3NGxGcmRqZ1BVUT09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
const testTlsCertBase64 = "Q2VydGlmaWNhdGU6CiAgICBEYXRhOgogICAgICAgIFZlcnNpb246IDEgKDB4MCkKICAgICAgICBTZXJpYWwgTnVtYmVyOgogICAgICAgICAgICAxODo2ZTpiMzo4NTo4NDplMDowMjoyMzpmNjo1YTpiMjoxMTplNDozOTozNzo0YTo0NDpkYzo5NjpiYQogICAgICAgIFNpZ25hdHVyZSBBbGdvcml0aG06IE5VTEwKICAgICAgICBJc3N1ZXI6IENOID0gc29tZS5jZXJ0LmRlCiAgICAgICAgVmFsaWRpdHkKICAgICAgICAgICAgTm90IEJlZm9yZTogTWFyIDIzIDEzOjM0OjE5IDIwMjEgR01UCiAgICAgICAgICAgIE5vdCBBZnRlciA6IE1hciAyMyAxMzozNDoxOSAyMDIyIEdNVAogICAgICAgIFN1YmplY3Q6IENOID0gc29tZS5jZXJ0LmRlCiAgICAgICAgU3ViamVjdCBQdWJsaWMgS2V5IEluZm86CiAgICAgICAgICAgIFB1YmxpYyBLZXkgQWxnb3JpdGhtOiByc2FFbmNyeXB0aW9uCiAgICAgICAgICAgICAgICBSU0EgUHVibGljLUtleTogKDIwNDggYml0KQogICAgICAgICAgICAgICAgTW9kdWx1czoKICAgICAgICAgICAgICAgICAgICAwMDpjZDoxMzo4Yjo5YTpmMjoyNzo3MzowMjoyZjowNTo2Mjo1Yjo1Njo0MToKICAgICAgICAgICAgICAgICAgICBlYzoxZjowMDoyNDpiNzo4ZDpiYjpiZDozYzoyZjowODoxZTowMzphZjo2YjoKICAgICAgICAgICAgICAgICAgICBlNTo3NDphMjpjNTozNTo2Zjo5NToxZjo2NzpmYTo4NzpiYTo0MjoxMDozYjoKICAgICAgICAgICAgICAgICAgICBkODo0MzpmYToxMjoyODo4ODoxYzozNDo1Njo3Mjo5YjoyZjpmZDpjMToxZToKICAgICAgICAgICAgICAgICAgICBmNjpkMDoxODpjNjpiNToxZDo5ODplMTozMDoxNTozMjozYjpjNTo5ZjoxMToKICAgICAgICAgICAgICAgICAgICA2Yzo4ZTo2Yjo4NjoxMTo5NTphMzoxMDoxNzpkNzo4MDpiYzo5NTo4Yjo5YjoKICAgICAgICAgICAgICAgICAgICA0YjphYjpkZjozYzo1ZDphYzozMjo3NzplMTpkNjowOTo1YTpjZDoyZTo1NzoKICAgICAgICAgICAgICAgICAgICAzMjpmNjoxNjpjZTo5ZTpjMzplNTo0ZTo4MTo4ZjowZDpkMzo1Zjo5YTpkYjoKICAgICAgICAgICAgICAgICAgICBiNDo1NTpjNzo4MDo3OToyYjo2NzoxNzpiMDplNjo5NTo5ODowMDoxMjo1ZjoKICAgICAgICAgICAgICAgICAgICAwOTowZDo1Yzo3ZDowZTozZDo3ODo2MzpiMToyYzozYjoxYjplOTpkMToyODoKICAgICAgICAgICAgICAgICAgICAzZDo2MjoyNDpkYjo5ZjpkNzphODpiMDpjNjozODozZToyMjoxYjpkMDo5ZDoKICAgICAgICAgICAgICAgICAgICBhNTo4NTpkZTo4OToyZTo1Zjo2NzozMTpjZjpkYjozMzo3Nzo5Yzo4ODphNDoKICAgICAgICAgICAgICAgICAgICAxOTpiMTo5ZDoxNTo0Mjo3ZDowZjpkNTo4OTpkNjo1Mjo4ZjozNjoxMTpmYToKICAgICAgICAgICAgICAgICAgICBkYTozYToyNjpmYjpiMTo1MzpmZDphMTo2Njo5OTpjMzo5MTo0NDo5MjpkMjoKICAgICAgICAgICAgICAgICAgICBjMjplODo0MDpiYzoxMDpjMjo3ZDpiODo2ZDpkMDphNTowNToxODo0ZDo3NToKICAgICAgICAgICAgICAgICAgICA1MTozZDpkZDpjNjo4NjpmMToxMzplNjowYTo1YTplNDoxYToyZTpiYjo3MjoKICAgICAgICAgICAgICAgICAgICBhMToyOTowMjoyYTozNDo1MjpkYjo3OTpkOTpjZDoxNjozNzo1ZjoyMDowYjoKICAgICAgICAgICAgICAgICAgICBlNTowNwogICAgICAgICAgICAgICAgRXhwb25lbnQ6IDY1NTM3ICgweDEwMDAxKQogICAgU2lnbmF0dXJlIEFsZ29yaXRobTogTlVMTAotLS0tLUJFR0lOIENFUlRJRklDQVRFLS0tLS0KTUlJQ3RUQ0NBWjBDRkRWY2Ywb3ZHOGZ4Ync2aTVjdWFDenhjRGpWWE1BMEdDU3FHU0liM0RRRUJDd1VBTUJjeApGVEFUQmdOVkJBTU1ESE52YldVdVkyVnlkQzVrWlRBZUZ3MHlNVEF6TWpNeE16TTBNVGxhRncweU1qQXpNak14Ck16TTBNVGxhTUJjeEZUQVRCZ05WQkFNTURITnZiV1V1WTJWeWRDNWtaVENDQVNJd0RRWUpLb1pJaHZjTkFRRUIKQlFBRGdnRVBBRENDQVFvQ2dnRUJBTTBUaTVyeUozTUNMd1ZpVzFaQjdCOEFKTGVOdTcwOEx3Z2VBNjlyNVhTaQp4VFZ2bFI5bitvZTZRaEE3MkVQNkVpaUlIRFJXY3Bzdi9jRWU5dEFZeHJVZG1PRXdGVEk3eFo4UmJJNXJoaEdWCm94QVgxNEM4bFl1YlM2dmZQRjJzTW5maDFnbGF6UzVYTXZZV3pwN0Q1VTZCanczVFg1cmJ0RlhIZ0hrclp4ZXcKNXBXWUFCSmZDUTFjZlE0OWVHT3hMRHNiNmRFb1BXSWsyNS9YcUxER09ENGlHOUNkcFlYZWlTNWZaekhQMnpOMwpuSWlrR2JHZEZVSjlEOVdKMWxLUE5oSDYyam9tKzdGVC9hRm1tY09SUkpMU3d1aEF2QkRDZmJodDBLVUZHRTExClVUM2R4b2J4RStZS1d1UWFMcnR5b1NrQ0tqUlMyM25aelJZM1h5QUw1UWNDQXdFQUFUQU5CZ2txaGtpRzl3MEIKQVFzRkFBT0NBUUVBTWJyRk1nZ0Z0LzZoR2l4OVcwOXQwS08xeFZHTTBlVVZ5S3c3VU1UWVc2YkdMaHgvZkZsbwpIOGdxWkw1bVM2cCt6alU5T1NabXBESGZnOVVhRERzT3pIQVA3SVhvWWtVYjRnMmIvSjVhdHA5QXJmcnliUDdZCjBua2drQ2ZIQXdrM1MydWhLNlR4UUJaM2JjcXJrUWJ4N3JGVGx5b2RzNGtKNHlRQjJFQzg1MkNYWVdGRWRMOHgKeWFHZXMyNGxkcUoyMmxKRldmRnhIZUNMTThQbHlhSzZsUHlnRy9vM01RSU9JSFJiMVNQSUh0VlU1NTlSZUJnMwp4c3VRdWQ3d25JTUwwRW5xWEVpSml0WXk1a3pBNEl4cmRkL2p0QW5GRkVBWjl1SDVqdzVzWWp6OVBZK3lOOHZ4Ck1YSVBJWGpQYWhGaHRWQjF4WU0vYWNOY1VHa1dFMFo4ZlE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="
const testTlsCertKeyBase64 = "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2QUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktZd2dnU2lBZ0VBQW9JQkFRRE5FNHVhOGlkekFpOEYKWWx0V1Fld2ZBQ1MzamJ1OVBDOElIZ092YStWMG9zVTFiNVVmWi9xSHVrSVFPOWhEK2hJb2lCdzBWbktiTC8zQgpIdmJRR01hMUhaamhNQlV5TzhXZkVXeU9hNFlSbGFNUUY5ZUF2SldMbTB1cjN6eGRyREozNGRZSldzMHVWekwyCkZzNmV3K1ZPZ1k4TjAxK2EyN1JWeDRCNUsyY1hzT2FWbUFBU1h3a05YSDBPUFhoanNTdzdHK25SS0QxaUpOdWYKMTZpd3hqZytJaHZRbmFXRjNva3VYMmN4ejlzemQ1eUlwQm14blJWQ2ZRL1ZpZFpTanpZUit0bzZKdnV4VS8yaApacG5Ea1VTUzBzTG9RTHdRd24yNGJkQ2xCUmhOZFZFOTNjYUc4UlBtQ2xya0dpNjdjcUVwQWlvMFV0dDUyYzBXCk4xOGdDK1VIQWdNQkFBRUNnZ0VBYWIycVBqcWVISzhEajhNblZWS29iVk9sbXY5NXpoazZKdlZTOFNDeEwzSysKUE05TUZPV0lTSFBCbkowKzVjNExqdHFmc0Z6aXV5SUR0WkJCc3dzVGFrL1loRVJHcWFBb1JkeTJITGxVWjd6QQpWNHZ6a20ycXJsRmtzenBuNWVUa0lPeFJjSUZoU29Pcnd6Zi9VZDJ3WHNwdStMSUVtZFN2SjR1MnNzT3VaSWZsCmU0NGkxOEQzTzNVbHFMaWh5VXFaMHk0b2xadjVzaWFtRXhqSEswbkE5cklKdTVxQktKdThTNWdKNnRqK3Y2YXAKdE5nUTJSNXhsNmk1QlpTT0wvbVdKdUNGK2RJajdYN2NydkRaTkZoZUxCK0VGVGtYZ29NV2ZMRlFLR29zYTBjagpFWi9uMEg3NStBTUg2WGJPS1U2QWEzaXVYNDVTbzZVSVVqai8rZFJnQVFLQmdRRDE5RExUYkRYUEU4QUZqNm1XClBJbDlGRFAwWCtEeThzRGtlL3RVdmxPZkRKZnZnVlNMNWY4eFlqalFPSEM2M3JBNTdOVC9nR3VkS0dLTnVCWUkKSHJBQnA5bUM1dXlPU0FDaktBZUNDbkNRY0lKNjBLV2FBY2ZkZE95UE5QV2UvbjM3WGg2TXJPazhuUzg3R25xYQpVYndBalpNZ0FZTXN5YUdFL2NjN05kb1Fnd0tCZ1FEVmMrbks2RkJpdXJyQzU0OXdydytwRzdHd0ZLTkx5RnVLCkdMZkQ2by9XeEEvODd3TFNWQVJsajJRQXBqS3dvUm43SUtUZ200QnJGZ2FBeTJqOU10Zml0czZUdTBtM21KUEMKMnJrMmVNQ1FPS0pYTlRRQm0ySXh1RXRYRDU1WXFUNER2ejFXRjFBQkRNSHBvZUFneHBaUHJ1RkVybmovdzYxSQpLV1hTRTBDcUxRS0JnQXowZStqZS9rYVdCN3REUWUrRDZNb0owbUxBMmh4eDVPOGtDSzBDQ1cyTFFFV0JUbTdBCkFwMGJTMXJNWGtPNWp4YTkvc29tZllTZHAvTkhDd0lLZThMYWtINXdvMjBySmIxeVVsTHJNZHFwMG5XZG45dG8KMUpvNW1teEFvZDlxRUVDNVNHcW9nUENNWnZ4NS9KTThVdWJFamtkVlRRK0MzMXNkOHV5UGZaajVBb0dBVkk3SQpyUkwrMVQyM3dvSk05b3pESFhEVklUWHJ3cGVxZTdoekErK2w3NlJYMlJFdUF2ZzVqYW9TS1pldE9QOTQ5VnpuCk0vc21Fa1gxYVl3ckdUTE5Cd2o0S05ubXlBNXZhcCtQQTU4dVdYTzJDK29Ob2gxVjl2QlZHRFlkdW0zQkhXYmkKKzNuY3ZhMjZHNzErdGowMVNuZXkwYXgwVG8zTDFXeGc0Nm13MGprQ2dZQUwvNHpPbXQ0c0ZzdEQ2OFY3WllhYgpiY3A5MExVNk56VnY5YzArbDh1SkhWWDNjdzc5QnJBS1hqYklxTDRiUnU1VXgzWmgvUUpISTVOa2p0a1BaeE40CjF3clhrd2hGNktZM1VYbzJSTEhCY0FiR0tPWC9IckdhczB2LzZZMG5Vc0tjQ3J0SFllU0pjdC9LVTQwREhaS3EKcVFZUlYyZ3dIZTRuYU01SXJsem54UT09Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0K"

func CreatePostgresTestProvider(t *testing.T) (*provider.IServiceProvider, config.RawConfig) {
	config := common_test.GetConfig(t)

	url, err := url.Parse(config.Kubernetes.Server)
	assert.Nil(t, err)

	var pgProvider provider.IServiceProvider = postgres.CreatePostgresProvider(
		url.Hostname(),
		config.Kubernetes.Server,
		config.Kubernetes.CertificateAuthority,
		config.Kubernetes.Operators.Postgres.PgoUrl,
		config.Kubernetes.Operators.Postgres.PgoVersion,
		config.Kubernetes.Operators.Postgres.PgoCaPath,
		config.Kubernetes.Operators.Postgres.PgoUsername,
		config.Kubernetes.Operators.Postgres.PgoPassword,
		config.ResourcesTemplatesPath,
		kubernetes.NginxInformation(config.Kubernetes.Nginx),
	)

	return &pgProvider, config
}

func Test_Postgres_Provider_Create_Panic_Template_Not_Found(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Postgres provider did not panic if the template could not be found")
		}
	}()

	pgProvider := postgres.CreatePostgresProvider(
		"Hostname",
		"Server",
		"CaPath",
		"NotExistingPath",
		"",
		"",
		"",
		"",
		"",
		kubernetes.NginxInformation{})

	_ = pgProvider
}

func Test_Postgres_Provider_GetAttributes(t *testing.T) {
	pgProviderPtr, _ := CreatePostgresTestProvider(t)
	pgProvider := *pgProviderPtr

	assert.NotEqual(t, "", pgProvider.GetServiceImage())
	assert.NotEqual(t, "", pgProvider.GetServiceDescription())
	assert.Equal(t, "Postgres", pgProvider.GetServiceType())

	testUser := unit_test.TestUser{
		KubernetesNamespace: "MyNamespace",
	}

	// Get json form data
	formDataObj1, err := pgProvider.GetJsonForm(testUser)
	assert.Nil(t, err)
	assert.NotNil(t, formDataObj1)
	formDataObj2, err := pgProvider.GetJsonForm(testUser)
	assert.Nil(t, err)
	assert.NotNil(t, formDataObj2)

	// Ensure they are not the same (because of the random name)
	formData1 := formDataObj1.(provider_dtos.FormQueryDto)
	formData2 := formDataObj2.(provider_dtos.FormQueryDto)

	assert.NotEqual(t,
		formData1.Properties.Common.Properties.ClusterName.Default,
		formData2.Properties.Common.Properties.ClusterName.Default)

	// Generate yaml from form values and ensure it sets the values from the form
	filledForm := provider_dtos.FormResponseDto{}
	filledForm.Common.ClusterName = "MyCluster"

	filledFormData, err := json.Marshal(filledForm)
	assert.Nil(t, err)
	yamlTemplate, err := pgProvider.GetYamlTemplate(testUser, filledFormData)
	assert.Nil(t, err)
	assert.NotNil(t, yamlTemplate)

	yamlObject := yamlTemplate.(provider_dtos.ProviderYamlTemplateDto)

	// Ensure values are set by form as expected
	expectedClusterName := "MyCluster"
	assert.Equal(t, expectedClusterName, yamlObject.Metadata.Annotations.CurrentPrimary)
	assert.Equal(t, expectedClusterName, yamlObject.Metadata.Labels.CrunchyPghaScope)
	assert.Equal(t, expectedClusterName, yamlObject.Metadata.Labels.DeploymentName)
	assert.Equal(t, expectedClusterName, yamlObject.Metadata.Labels.Name)
	assert.Equal(t, expectedClusterName, yamlObject.Metadata.Labels.PgCluster)
	assert.Equal(t, expectedClusterName, yamlObject.Metadata.Name)
	assert.Equal(t, expectedClusterName, yamlObject.Spec.PrimaryStorage.Name)
	assert.Equal(t, expectedClusterName, yamlObject.Spec.Clustername)
	assert.Equal(t, expectedClusterName, yamlObject.Spec.Database)
	assert.Equal(t, expectedClusterName, yamlObject.Spec.Name)

	expectedNamespace := "MyNamespace"
	assert.Equal(t, expectedNamespace, yamlObject.Metadata.Namespace)
	assert.Equal(t, expectedNamespace, yamlObject.Spec.Namespace)
}

func Test_Postgres_Provider_End2End_Tls_From_File(t *testing.T) {
	pgProviderPtr, config := CreatePostgresTestProvider(t)
	user := config.Users[0]

	filledForm := provider_dtos.FormResponseDto{}
	filledForm.Common.InClusterPort = 5432
	filledForm.Common.Username = "testuser"
	filledForm.Common.ClusterStorageSize = 1

	// Tls from file
	filledForm.TLS.UseTLS = true
	filledForm.Common.ClusterName = "pg-test1-cluster"
	filledForm.TLS.TLSMode = "TlsFromFile"
	filledForm.TLS.TLSModeFromFile.CaCertBase64 = testCaBase64
	filledForm.TLS.TLSModeFromFile.TlsCertificateBase64 = testTlsCertBase64
	filledForm.TLS.TLSModeFromFile.TlsPrivateKeyBase64 = testTlsCertKeyBase64

	Postgres_Provider_End2End(t, pgProviderPtr, user, filledForm)
}

func Test_Postgres_Provider_End2End_Tls_From_Secret(t *testing.T) {
	pgProviderPtr, config := CreatePostgresTestProvider(t)
	user := config.Users[0]

	api, err := kubernetes.GenerateK8sApiFromToken(config.Kubernetes.Server, config.Kubernetes.CertificateAuthority, user.GetKubernetesAccessToken())
	assert.Nil(t, err)

	caCrt, err := base64.StdEncoding.DecodeString(testCaBase64)
	assert.Nil(t, err)

	_, err = api.CreateTlsSecretWithoutOwner(
		"pg-test2-cluster-ca",
		user.GetKubernetesNamespace(),
		map[string][]byte{
			"ca.crt": caCrt,
		},
	)
	assert.Nil(t, err)

	tlsCrt, err := base64.StdEncoding.DecodeString(testTlsCertBase64)
	assert.Nil(t, err)

	tlsKey, err := base64.StdEncoding.DecodeString(testTlsCertKeyBase64)
	assert.Nil(t, err)

	_, err = api.CreateTlsSecretWithoutOwner(
		"pg-test2-cluster-tls-keypair",
		user.GetKubernetesNamespace(),
		map[string][]byte{
			"tls.crt": tlsCrt,
			"tls.key": tlsKey,
		},
	)
	assert.Nil(t, err)

	time.Sleep(20 * time.Second)

	// Tls fromm secret
	filledForm := provider_dtos.FormResponseDto{}
	filledForm.Common.InClusterPort = 5432
	filledForm.Common.Username = "testuser"
	filledForm.Common.ClusterStorageSize = 1
	filledForm.Common.ClusterName = "pg-test2-cluster"

	filledForm.TLS.UseTLS = true
	filledForm.TLS.TLSMode = "TlsFromSecret"
	filledForm.TLS.TLSModeFromSecret.TLSSecret = "pg-test2-cluster-tls-keypair"
	filledForm.TLS.TLSModeFromSecret.CaSecret = "pg-test2-cluster-ca"

	Postgres_Provider_End2End(t, pgProviderPtr, user, filledForm)

	err = api.DeleteSecret(user.GetKubernetesNamespace(), "pg-test2-cluster-ca")
	assert.Nil(t, err)

	err = api.DeleteSecret(user.GetKubernetesNamespace(), "pg-test2-cluster-tls-keypair")
	assert.Nil(t, err)
}

func Test_Postgres_Provider_End2End_NoTls(t *testing.T) {
	pgProviderPtr, config := CreatePostgresTestProvider(t)
	user := config.Users[0]

	filledForm := provider_dtos.FormResponseDto{}
	filledForm.Common.ClusterName = "pg-test3-cluster"
	filledForm.Common.InClusterPort = 5432
	filledForm.Common.Username = "testuser"
	filledForm.Common.ClusterStorageSize = 1
	filledForm.TLS.UseTLS = false

	Postgres_Provider_End2End(t, pgProviderPtr, user, filledForm)
}

func Postgres_Provider_End2End(t *testing.T, pgProviderPtr *provider.IServiceProvider, user common.IKubernetesAuthInformation, filledForm provider_dtos.FormResponseDto) {
	service1Ptr := common_test.CommonProviderStart(t, pgProviderPtr, user, filledForm, 4)


	// --- user ---
	// show = 1
	actionPtr, err := common_test.GetAction(service1Ptr, "User", "cmd_pg_show_users")
	assert.Nil(t, err)
	action := *actionPtr
	result, err := action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.Nil(t, err)
	users := result.(map[string]string)
	assert.Equal(t, 1, len(users))

	// Add 1
	actionPtr, err = common_test.GetAction(service1Ptr, "User", "cmd_pg_add_user")
	assert.Nil(t, err)
	action = *actionPtr
	addUserDto := *(action.GetJsonFormResultPlaceholder().(*action_dtos.AddUserDto))
	addUserDto.Password = "testpswd"
	addUserDto.Username = "testuser1"
	_, err = action.GetActionExecuteCallback()(&addUserDto)
	assert.Nil(t, err)
	time.Sleep(5 * time.Second)

	// show = 2
	actionPtr, err = common_test.GetAction(service1Ptr, "User", "cmd_pg_show_users")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.Nil(t, err)
	users = result.(map[string]string)
	assert.Equal(t, 2, len(users))

	// Delete 1
	actionPtr, err = common_test.GetAction(service1Ptr, "User", "cmd_pg_remove_user")
	assert.Nil(t, err)
	action = *actionPtr
	deleteUserDto := *(action.GetJsonFormResultPlaceholder().(*action_dtos.DeleteUserDto))
	deleteUserDto.Username = "testuser1"
	_, err = action.GetActionExecuteCallback()(&deleteUserDto)
	assert.Nil(t, err)
	time.Sleep(5 * time.Second)

	// show = 1
	actionPtr, err = common_test.GetAction(service1Ptr, "User", "cmd_pg_show_users")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.Nil(t, err)
	users = result.(map[string]string)
	assert.Equal(t, 1, len(users))

	// --- Exposure ---
	// Check if toggle is correct
	toggleActionPtr, err := common_test.GetToggleAction(service1Ptr, "Security", "cmd_pg_expose_toggle")
	assert.Nil(t, err)
	toggleAction := *toggleActionPtr
	isSet, err := toggleAction.Get()
	assert.Nil(t, err)
	assert.False(t, isSet) // Not exposed

	// Check expose details
	actionPtr, err = common_test.GetAction(service1Ptr, "Security", "cmd_pg_get_expose_info")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.Nil(t, err)
	clusterExposeInformation := result.(*action_dtos.ClusterExposeResponseDto)
	assert.Equal(t, 0, clusterExposeInformation.Port)

	// Expose it
	toggleActionPtr, err = common_test.GetToggleAction(service1Ptr, "Security", "cmd_pg_expose_toggle")
	assert.Nil(t, err)
	toggleAction = *toggleActionPtr
	result, err = toggleAction.Set()
	assert.Nil(t, err)
	assert.Nil(t, result)
	time.Sleep(5 * time.Second)

	// Check if toggle is correct
	toggleActionPtr, err = common_test.GetToggleAction(service1Ptr, "Security", "cmd_pg_expose_toggle")
	assert.Nil(t, err)
	toggleAction = *toggleActionPtr
	isSet, err = toggleAction.Get()
	assert.Nil(t, err)
	assert.True(t, isSet) // exposed

	// Check expose details
	actionPtr, err = common_test.GetAction(service1Ptr, "Security", "cmd_pg_get_expose_info")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.Nil(t, err)
	clusterExposeInformation = result.(*action_dtos.ClusterExposeResponseDto)
	assert.True(t, clusterExposeInformation.Port > 0)

	// Hide it again
	toggleActionPtr, err = common_test.GetToggleAction(service1Ptr, "Security", "cmd_pg_expose_toggle")
	assert.Nil(t, err)
	toggleAction = *toggleActionPtr
	result, err = toggleAction.Unset()
	assert.Nil(t, err)
	time.Sleep(5 * time.Second)

	// Check again if it is hidden
	// Check if toggle is correct
	toggleActionPtr, err = common_test.GetToggleAction(service1Ptr, "Security", "cmd_pg_expose_toggle")
	assert.Nil(t, err)
	toggleAction = *toggleActionPtr
	isSet, err = toggleAction.Get()
	assert.Nil(t, err)
	assert.False(t, isSet) // Not exposed

	// Check expose details
	actionPtr, err = common_test.GetAction(service1Ptr, "Security", "cmd_pg_get_expose_info")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.Nil(t, err)
	clusterExposeInformation = result.(*action_dtos.ClusterExposeResponseDto)
	assert.Equal(t, 0, clusterExposeInformation.Port)

	// --- Scale ---
	actionPtr, err = common_test.GetAction(service1Ptr, "Features", "cmd_pg_scale")
	assert.Nil(t, err)
	action = *actionPtr
	clusterScale := *(action.GetJsonFormResultPlaceholder().(*action_dtos.ClusterScaleDto))
	assert.Equal(t, 0, clusterScale.NumberOfReplicas)
	// Try setting the same number of replicas as we have
	clusterScale.NumberOfReplicas = 0
	result, err = action.GetActionExecuteCallback()(&clusterScale)
	assert.NotNil(t, err) // Should create an error
	assert.Nil(t, result)
	// Try setting a negative number of replicas
	clusterScale.NumberOfReplicas = -1
	result, err = action.GetActionExecuteCallback()(&clusterScale)
	assert.NotNil(t, err) // Should create an error
	assert.Nil(t, result)
	// Increment the number of replicas
	clusterScale.NumberOfReplicas = 2
	result, err = action.GetActionExecuteCallback()(&clusterScale)
	assert.Nil(t, err)
	assert.Nil(t, result)
	time.Sleep(5 * time.Second)
	// Ensure we have 2 replicas now
	assert.Nil(t, err)
	actionPtr, err = common_test.GetAction(service1Ptr, "Features", "cmd_pg_scale")
	assert.Nil(t, err)
	action = *actionPtr
	clusterScale = *(action.GetJsonFormResultPlaceholder().(*action_dtos.ClusterScaleDto))
	assert.Equal(t, 2, clusterScale.NumberOfReplicas)
	// Decrement the number of replicas
	clusterScale.NumberOfReplicas = 1
	result, err = action.GetActionExecuteCallback()(&clusterScale)
	assert.Nil(t, err)
	assert.Nil(t, result)
	time.Sleep(5 * time.Second)
	// Ensure we have only 1 replica now
	assert.Nil(t, err)
	actionPtr, err = common_test.GetAction(service1Ptr, "Features", "cmd_pg_scale")
	assert.Nil(t, err)
	action = *actionPtr
	clusterScale = *(action.GetJsonFormResultPlaceholder().(*action_dtos.ClusterScaleDto))
	assert.Equal(t, 1, clusterScale.NumberOfReplicas)

	// Shut down all services of the provider
	common_test.CommonProviderStop(t, pgProviderPtr, user)
}
