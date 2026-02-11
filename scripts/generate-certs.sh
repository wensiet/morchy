#!/bin/bash

set -e

# Default values
CERT_DIR="./certs"
VALIDITY_YEARS=2
CONTROL_PLANE_NODE_ID="controlplane-1"
AGENT_NODE_ID=""
GENERATE_AGENT=false

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Parse arguments
if [[ $# -gt 0 && "$1" == "--agent" ]]; then
	GENERATE_AGENT=true
	if [[ $# -lt 2 ]]; then
		echo -e "${YELLOW}Usage: $0 --agent <agent-node-id>${NC}"
		exit 1
	fi
	AGENT_NODE_ID="$2"
	VALIDITY_YEARS=${3:-2}
else
	VALIDITY_YEARS=${1:-2}
fi

echo -e "${BLUE}=== mTLS Certificate Generation Script ===${NC}"
echo "Certificate directory: $CERT_DIR"
echo "Validity: $VALIDITY_YEARS years"

if [[ "$GENERATE_AGENT" == true ]]; then
	echo "Mode: Generate CA + Server + Agent certificates"
	echo "Agent Node ID: $AGENT_NODE_ID"
else
	echo "Mode: Generate CA + Control Plane certificate only"
fi
echo ""

# Create certificate directory
mkdir -p "$CERT_DIR"

# Generate CA private key
echo -e "${BLUE}Generating CA private key...${NC}"
openssl genrsa -out "$CERT_DIR/ca-key.pem" 2048 2>/dev/null

# Generate CA certificate
echo -e "${BLUE}Generating CA certificate...${NC}"
openssl req -new -x509 -key "$CERT_DIR/ca-key.pem" \
  -out "$CERT_DIR/ca-cert.pem" \
  -days $((VALIDITY_YEARS * 365)) \
  -subj "/C=US/ST=State/L=City/O=Morchy/CN=Morchy-CA" 2>/dev/null

echo -e "${GREEN}✓ CA certificate created${NC}"
echo ""

# Generate Control Plane server certificate
echo -e "${BLUE}Generating Control Plane server certificate...${NC}"

# Generate private key
openssl genrsa -out "$CERT_DIR/controlplane-key.pem" 2048 2>/dev/null

# Create certificate signing request
openssl req -new \
  -key "$CERT_DIR/controlplane-key.pem" \
  -out "$CERT_DIR/controlplane.csr" \
  -subj "/C=US/ST=State/L=City/O=Morchy/CN=$CONTROL_PLANE_NODE_ID" 2>/dev/null

# Create extensions file for server certificate
cat > "$CERT_DIR/controlplane-ext.conf" << EOF
subjectAltName=DNS:localhost,DNS:127.0.0.1,IP:127.0.0.1
keyUsage=digitalSignature,keyEncipherment
extendedKeyUsage=serverAuth
EOF

# Sign the certificate with CA
openssl x509 -req \
  -in "$CERT_DIR/controlplane.csr" \
  -CA "$CERT_DIR/ca-cert.pem" \
  -CAkey "$CERT_DIR/ca-key.pem" \
  -CAcreateserial \
  -out "$CERT_DIR/controlplane-cert.pem" \
  -days $((VALIDITY_YEARS * 365)) \
  -extfile "$CERT_DIR/controlplane-ext.conf" 2>/dev/null

echo -e "${GREEN}✓ Control Plane server certificate created${NC}"
echo ""

if [[ "$GENERATE_AGENT" == true ]]; then
	# Generate Agent client certificate
	echo -e "${BLUE}Generating Agent client certificate (Node ID: $AGENT_NODE_ID)...${NC}"

	# Determine file names for agent certificate (support multiple agents)
	AGENT_CERT_FILE="$CERT_DIR/agent-${AGENT_NODE_ID}-cert.pem"
	AGENT_KEY_FILE="$CERT_DIR/agent-${AGENT_NODE_ID}-key.pem"
	AGENT_CSR_FILE="$CERT_DIR/agent-${AGENT_NODE_ID}.csr"

	# Generate private key
	openssl genrsa -out "$AGENT_KEY_FILE" 2048 2>/dev/null

	# Create certificate signing request with node ID in CommonName
	openssl req -new \
	  -key "$AGENT_KEY_FILE" \
	  -out "$AGENT_CSR_FILE" \
	  -subj "/C=US/ST=State/L=City/O=Morchy/CN=$AGENT_NODE_ID" 2>/dev/null

	# Create extensions file for client certificate
	cat > "$CERT_DIR/agent-ext.conf" << EOF
keyUsage=digitalSignature,keyEncipherment
extendedKeyUsage=clientAuth
EOF

	# Sign the certificate with CA
	openssl x509 -req \
	  -in "$AGENT_CSR_FILE" \
	  -CA "$CERT_DIR/ca-cert.pem" \
	  -CAkey "$CERT_DIR/ca-key.pem" \
	  -CAcreateserial \
	  -out "$AGENT_CERT_FILE" \
	  -days $((VALIDITY_YEARS * 365)) \
	  -extfile "$CERT_DIR/agent-ext.conf" 2>/dev/null

	echo -e "${GREEN}✓ Agent client certificate created: $(basename $AGENT_CERT_FILE)${NC}"
	echo ""

	rm -f "$AGENT_CSR_FILE" "$CERT_DIR/agent-ext.conf"
fi

# Clean up temporary files
rm -f "$CERT_DIR/controlplane.csr" "$CERT_DIR/controlplane-ext.conf" "$CERT_DIR/ca-key.srl"

# Verify certificates
echo -e "${BLUE}Verifying certificates...${NC}"
openssl verify -CAfile "$CERT_DIR/ca-cert.pem" "$CERT_DIR/controlplane-cert.pem" 2>/dev/null && \
  echo -e "${GREEN}✓ Control Plane certificate verified${NC}" || \
  echo -e "✗ Control Plane certificate verification failed"

if [[ "$GENERATE_AGENT" == true ]]; then
	openssl verify -CAfile "$CERT_DIR/ca-cert.pem" "$AGENT_CERT_FILE" 2>/dev/null && \
	  echo -e "${GREEN}✓ Agent certificate verified: $AGENT_NODE_ID${NC}" || \
	  echo -e "✗ Agent certificate verification failed"
fi

echo ""

# Display certificate information
echo -e "${BLUE}=== Certificate Information ===${NC}"
echo ""
echo "CA Certificate:"
openssl x509 -in "$CERT_DIR/ca-cert.pem" -text -noout | grep -A2 "Subject:" | head -2

echo ""
echo "Control Plane Server Certificate:"
openssl x509 -in "$CERT_DIR/controlplane-cert.pem" -text -noout | grep -A2 "Subject:" | head -2

if [[ "$GENERATE_AGENT" == true ]]; then
	echo ""
	echo "Agent Client Certificate (Node ID: $AGENT_NODE_ID):"
	openssl x509 -in "$AGENT_CERT_FILE" -text -noout | grep -A2 "Subject:" | head -2
fi

echo ""
echo -e "${GREEN}✓ Certificate generation completed!${NC}"
