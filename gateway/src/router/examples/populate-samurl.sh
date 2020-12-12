#!/bin/bash
curl -ik -X DELETE https://localhost:8443/services/v1/samurl/bicycleurl -u samurl:ddd
curl -d "@table.json" -X POST https://localhost:8443/services/v1/samurl/bicycleurl.csv -ik -u samurl:ddd
curl -d "@tablecols.en.json" -X POST https://localhost:8443/services/v1/samurl/bicycleurl/colnames/en -ik -u samurl:ddd
curl -d "@tablecols.it.json" -X POST https://localhost:8443/services/v1/samurl/bicycleurl/colnames/it -ik -u samurl:ddd
curl -d "@tablevalues.json" -X POST https://localhost:8443/services/v1/samurl/bicycleurl/values -ik -u samurl:ddd
curl -d "@tablevalues.2.json" -X POST https://localhost:8443/services/v1/samurl/bicycleurl/values -ik -u samurl:ddd
curl -d "@table.enabled.json" -X PUT https://localhost:8443/services/v1/samurl/bicycleurl -ik -u samurl:ddd
curl -ik -X GET https://localhost:8443/services/v1/samurl/bicycleurl
curl -d "@table.1.json" -X POST https://localhost:8443/services/v1/samurl/otherservice.csv -ik -u samurl:ddd
curl -d "@table.2.json" -X POST https://localhost:8443/services/v1/samurl/service2.csv -ik -u samurl:ddd
curl -d "@table.3.json" -X POST https://localhost:8443/services/v1/samurl/service3.csv -ik -u samurl:ddd
curl -d "@table.4.json" -X POST https://localhost:8443/services/v1/samurl/service4.csv -ik -u samurl:ddd