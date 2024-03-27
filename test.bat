@echo off

go run . -f .\examples\files\hgf.txt -r .\examples\rules.yaml
go run . -f .\examples\files\amazon_web_services.txt -r .\examples\rules.yaml
go run . -f .\examples\files\google_cloud.txt -r .\examples\rules.yaml
go run . -f .\examples\files\paidonresults.txt -r .\examples\rules.yaml
go run . -f .\examples\files\test_data.txt -r .\examples\test_data_rules\rule1.yaml
go run . -f .\examples\files\test_data.txt -r .\examples\test_data_rules\rule2.yaml
go run . -f .\examples\files\test_data.txt -r .\examples\test_data_rules\rule3.yaml
go run . -f .\examples\files\test_data.txt -r .\examples\test_data_rules\rule4.yaml
go run . -f .\examples\files\test_data.txt -r .\examples\test_data_rules\rule5.yaml
go run . -f .\examples\files\test_data.txt -r .\examples\test_data_rules\rule6.yaml
go run . -f .\examples\files\test_data.txt -r .\examples\test_data_rules\rule7.yaml
go run . -f .\examples\files\test_data.txt -r .\examples\test_data_rules\rule8.yaml
go run . -f .\examples\files\test_data.txt -r .\examples\test_data_rules\rule9.yaml

pause