go run ./cmd/hdarrrr -tonemapper reinhard05 -low examples/landscape_low.jpeg -mid examples/landscape_mid.jpeg -high examples/landscape_high.jpeg -output examples/landscape_hdr-reinhard05.jpg



go run ./cmd/hdarrrr -tonemapper drago03 -low examples/landscape_low.jpeg -mid examples/landscape_mid.jpeg -high examples/landscape_high.jpeg -output examples/landscape_hdr-drago03.jpg


open examples/landscape_hdr-reinhard05.jpg
open examples/landscape_hdr-drago03.jpg
