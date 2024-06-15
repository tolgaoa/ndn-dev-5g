#!/bin/bash

# Check if the number of requests was passed as an argument
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <number_of_requests>"
    exit 1
fi

# Total number of requests to distribute over a 1-minute window
TOTAL_REQUESTS=$1

DURATION_MINUTES=10

# Output files
outputFile="request_times_with_packet_loss.csv"

# Get pod names
namespace="oai"
amfpod=$(kubectl get pods -n $namespace | grep amf | awk '{print $1}')

# Function to initialize output files
initialize_output_file() {
    if [ ! -f "$1" ]; then
        echo "Timestamp,HTTP Version,Packet Loss Rate,Response Time (ms),Success" > "$1"
    fi
}

# Initialize output files
initialize_output_file "$outputFile"

# Functions to make specific requests and log the response times
make_request() {
    local http_version=$1
    local packet_loss_rate=$2
    local endpoint=$3
    local method=$4
    local data=$5
    local timestamp=$(date +%Y-%m-%dT%H:%M:%S)

    # Apply packet loss
    kubectl exec $amfpod -n $namespace -c proxy -- tc qdisc add dev eth0 root netem loss $packet_loss_rate

    # Make the request and measure response time
    local response=$(kubectl exec $amfpod -n $namespace -- curl -X $method "$endpoint" -d "$data" -H "Content-Type: application/json" -w '%{time_total}' -o /dev/null -s)
    local success=$?

    # Remove packet loss
    kubectl exec $amfpod -n $namespace -c proxy -- tc qdisc del dev eth0 root

    # Calculate response time in ms
    local time_ms=$(echo "$response" | awk '{print $1*1000}')

    # Log the result
    echo "$timestamp,$http_version,$packet_loss_rate,$time_ms,$success" >> "$outputFile"
}

# Function to execute requests for a minute
execute_requests_for_a_minute() {
    local interval=$(echo "scale=2; 60 / $TOTAL_REQUESTS" | bc) # Interval between requests to fit into 1 minute
    local http_version=$1
    local packet_loss_rate=$2

    for (( i=0; i<TOTAL_REQUESTS; i++ )); do
        # Determine which request to make based on modulo of i
        case $((i % 4)) in
            0) make_request $http_version $packet_loss_rate "http://oai-ausf10-svc:80/nausf-auth/v1/ue-authentications" "POST" '{"servingNetworkName":"5G:mnc095.mcc208.3gppnetwork.org","supiOrSuci":"208950000000010"}' & ;;
            1) make_request $http_version $packet_loss_rate "http://oai-nrf10-svc:80/nnrf-disc/v1/nf-instances?target-nf-type=SMF&requester-nf-type=AMF" "GET" "" & ;;
            2) make_request $http_version $packet_loss_rate "http://oai-smf10-svc:80/nsmf-pdusession/v1/sm-contexts" "POST" '{"anType":"3GPP_ACCESS","dnn":"oai","gpsi":"msisdn-200000000001","n1MessageContainer":{"n1MessageClass":"SM","n1MessageContent":{"contentId":"n1SmMsg"}},"pduSessionId":1,"pei":"imei-200000000000001","requestType":"INITIAL_REQUEST","sNssai":{"sd":"123","sst":210},"servingNetwork":{"mcc":"208","mnc":"95"},"servingNfId":"servingNfId","smContextStatusUri":"http://10.42.0.35:80/nsmf-pdusession/callback/imsi-208950000000010/1","supi":"imsi-208950000000010"}' & ;;
            3) make_request $http_version $packet_loss_rate "http://oai-smf10-svc:80/nsmf-pdusession/v1/sm-contexts/1/modify" "POST" '{"n2SmInfo":{"contentId":"n2msg"},"n2SmInfoType":"PDU_RES_SETUP_RSP"}' & ;;
        esac

        sleep "$interval" # Control the rate of requests
    done
    wait # Wait for all background processes to complete
}

# Function to execute requests for the specified duration
execute_requests_for_duration() {
    local http_version=$1
    local packet_loss_rate=$2
    for (( j=0; j<DURATION_MINUTES; j++ )); do
        execute_requests_for_a_minute $http_version $packet_loss_rate
        echo "Cycle $((j+1)) of $DURATION_MINUTES completed."
        sleep 1 # Small delay before starting the next cycle
    done
}

# Define the HTTP versions and packet loss rates to test
http_versions=("HTTP1" "HTTP2" "HTTP3")
packet_loss_rates=("0%" "1%" "5%" "10%")

# Execute the requests for each combination of HTTP version and packet loss rate
for http_version in "${http_versions[@]}"; do
    for packet_loss_rate in "${packet_loss_rates[@]}"; do
        echo "Testing $http_version with $packet_loss_rate packet loss"
        execute_requests_for_duration $http_version $packet_loss_rate
    done
done

echo "All requests executed for $DURATION_MINUTES minutes."

