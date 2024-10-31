go run ./cmd/hdarrrr -method exposure-fusion -low examples/landscape_low.jpeg -mid examples/landscape_mid.jpeg -high examples/landscape_high.jpeg -output examples/landscape_hdr-exposure-fusion.jpg

go run ./cmd/hdarrrr -method tone-mapping -low examples/landscape_low.jpeg -mid examples/landscape_mid.jpeg -high examples/landscape_high.jpeg -output examples/landscape_hdr-tone-mapping.jpg
