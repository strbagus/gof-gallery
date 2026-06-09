#!/bin/bash

# ==============================================================================
# CONFIGURATION
# ==============================================================================
SOURCE_DIR="/mnt/usbstorage/NAS/Photos/Trah_Projo_2509"
EVENT_SLUG="trah-projo-2509"

# S3 Targets
# S3_ORIGINALS="ls3rpi/gallery/originals/$EVENT_SLUG"
# S3_PREVIEWS="ls3rpi/gallery/previews/$EVENT_SLUG"

# Homelab API & Telegram Endpoints
DB_API_ENDPOINT="http://192.168.99.108:3000/api/v1/photos"
TELEGRAM_ENDPOINT="http://192.168.99.108:8055/sendMessage"
BOT_CODE="123456"

# Local Workspace for temporary resizing
# TMP_DIR="/tmp/photo_processing"
# mkdir -p "$TMP_DIR"

# ==============================================================================
# INITIAL SANITY CHECKS
# ==============================================================================
if [ ! -d "$SOURCE_DIR" ]; then
    echo "❌ Error: Source directory $SOURCE_DIR does not exist!"
    exit 1
fi

echo "🚀 Starting Photography Processing Pipeline [Pattern: JSONB]..."
echo "--------------------------------------------------"

TOTAL_FILES=$(ls -1 "$SOURCE_DIR" | wc -l)
COUNTER=0

# ==============================================================================
# ITERATION LOOP
# ==============================================================================
for file in "$SOURCE_DIR"/*; do
    if [ -f "$file" ]; then
        COUNTER=$((COUNTER + 1))
        FILENAME=$(basename "$file")
        BASE_NAME="${FILENAME%.*}" # Strip file extension
        
        echo "📸 [$COUNTER/$TOTAL_FILES] Processing: $FILENAME"
        
        # ----------------------------------------------------------------------
        # STEP 1: Upload untouched original file
        # ----------------------------------------------------------------------
        # mc cp "$file" "$S3_ORIGINALS/$FILENAME" > /dev/null 2>&1
        # if [ $? -ne 0 ]; then
            # echo "❌ Failed to upload original file: $FILENAME. Skipping."
            # continue
        # fi

        # ----------------------------------------------------------------------
        # STEP 2: Generate WebP Responsive Formats using vips
        # ----------------------------------------------------------------------
        # vips thumbnail "$file" "$TMP_DIR/lg_${BASE_NAME}.webp" 1600
        # vips thumbnail "$file" "$TMP_DIR/md_${BASE_NAME}.webp" 400
        # vips thumbnail "$file" "$TMP_DIR/sm_${BASE_NAME}.webp" 200

        # ----------------------------------------------------------------------
        # STEP 3: Upload WebP Previews to S3
        # ----------------------------------------------------------------------
        # mc cp "$TMP_DIR/lg_${BASE_NAME}.webp" "$S3_PREVIEWS/lg_${BASE_NAME}.webp" > /dev/null 2>&1
        # mc cp "$TMP_DIR/md_${BASE_NAME}.webp" "$S3_PREVIEWS/md_${BASE_NAME}.webp" > /dev/null 2>&1
        # mc cp "$TMP_DIR/sm_${BASE_NAME}.webp" "$S3_PREVIEWS/sm_${BASE_NAME}.webp" > /dev/null 2>&1

        # ----------------------------------------------------------------------
        # STEP 4: Build Nested JSON Data Payload
        # ----------------------------------------------------------------------
        # ORIGINAL_URL="/gallery/originals/$EVENT_SLUG/$FILENAME"
        # LG_URL="/gallery/previews/$EVENT_SLUG/lg_${BASE_NAME}.webp"
        # MD_URL="/gallery/previews/$EVENT_SLUG/md_${BASE_NAME}.webp"
        # SM_URL="/gallery/previews/$EVENT_SLUG/sm_${BASE_NAME}.webp"

        # ----------------------------------------------------------------------
        # STEP 5: Sync to Go API Endpoint
        # ----------------------------------------------------------------------
        curl -s -X POST "$DB_API_ENDPOINT" \
          -H "Content-Type: application/json" \
          -d "{
            \"event_slug\": \"$EVENT_SLUG\",
            \"filename\": \"$FILENAME\",
          }" > /dev/null

        # ----------------------------------------------------------------------
        # STEP 6: Telegram Notification Alert
        # ----------------------------------------------------------------------
        # TELEGRAM_MSG="Processed [$COUNTER/$TOTAL_FILES]: $FILENAME (Original + JSONB Responsive Assets Synchronized)"
        # ENCODED_MSG=$(echo -v "$TELEGRAM_MSG" | curl -s -o /dev/null -w "%{url_effective}" --get --data-urlencode "message=$TELEGRAM_MSG" "" | cut -c 3-)
        # curl -s "$TELEGRAM_ENDPOINT?message=${ENCODED_MSG}&code=${BOT_CODE}" > /dev/null

        # ----------------------------------------------------------------------
        # STEP 7: Workspace Asset Eviction
        # ----------------------------------------------------------------------
        # rm -f "$TMP_DIR/lg_${BASE_NAME}.webp" "$TMP_DIR/md_${BASE_NAME}.webp" "$TMP_DIR/sm_${BASE_NAME}.webp"
        
        echo "✅ Finished loop execution for $FILENAME"
        echo "--------------------------------------------------"
    fi
done

# Clean up workspace completely
rm -rf "$TMP_DIR"

# Final batch summary transmission
FINAL_STATUS="🏁 Import Complete! $TOTAL_FILES images processed using JSONB data patterns and synced to the core database."
ENCODED_FINAL=$(echo -v "$FINAL_STATUS" | curl -s -o /dev/null -w "%{url_effective}" --get --data-urlencode "message=$FINAL_STATUS" "" | cut -c 3-)
curl -s "$TELEGRAM_ENDPOINT?message=${ENCODED_FINAL}&code=${BOT_CODE}" > /dev/null

echo "🏁 Pipeline complete!"
