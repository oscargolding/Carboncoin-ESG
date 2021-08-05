#!/bin/bash

base="localhost:3000/api"

# Register the users on the application

echo "Making the requests for users"
create="${base}/admin/auth/register"
JsonRequest () {
    printf '{"email":"%s","firm":"%s","password":"%s"}\n' "$1" "$2" "$3"
}
firstRequest=$(JsonRequest "oscarIndustry" "oscar" "oscar")
secondRequest=$(JsonRequest "BHPPetrol" "oscar" "oscar")
thirdRequest=$(JsonRequest "riotinto@gmail.com" "oscar" "oscar")

# Send the request
first_token=$(curl -X POST -H "Content-Type: application/json" \
    -d "${firstRequest}" \
    ${create} | jq -r '.token')

echo "Got first token ${first_token}"

second_token=$(curl -X POST -H "Content-Type: application/json" \
    -d "${secondRequest}" \
    ${create} | jq -r '.token')

echo "Got second token ${second_token}"

third_token=$(curl -X POST -H "Content-Type: application/json" \
    -d "${thirdRequest}" \
    ${create} | jq -r '.token')

echo "Got third token ${third_token}"
echo "Finished the requests for user creation"

echo "Reporting production of hydrogen"
rm -rf wallet
node certifier.js
echo "Finishing production of hydrogen"

echo "Create the offers on the market"

JsonRequestOffer () {
    printf '{"amount":"%s","tokens":"%s"}\n' "$1" "$2"
}

create_offer="${base}/offer/create"

CreateOffers () {
    first=("20" "30" "40" "50")
    second=("5" "10" "15" "20")
    for i in "${!first[@]}"; do
        printf "%s amount with %s tokens\n" "${first[i]}" "${second[i]}"
        echo "Making offer call with token $1"
        request=$(JsonRequestOffer "${first[i]}" "${second[i]}")
        auth="Authorization: Bearer $1"
        curl -X POST -H "Content-Type: application/json" -H "$auth" \
            -d "${request}" ${create_offer}
        printf "\nFinished Call"
    done
}

CreateOffers "$first_token"
CreateOffers "$second_token"
CreateOffers "$third_token"