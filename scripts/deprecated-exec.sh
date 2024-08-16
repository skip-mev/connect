#!/busybox/sh

echo "
     _
    | |
 ___| |_ ___  _ __
/ __| __/ _ \| '_ \
\__ \ || (_) | |_) |
|___/\__\___/| .__/
             | |
             |_|

WARNING: This container is deprecated and will be removed in future releases.
Please migrate to the corresponding Connect image."

# Execute the binary passed as arguments
exec "$@"