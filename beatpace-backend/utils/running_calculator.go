package utils

// StrideLength calculates the approximate stride length in meters based on height and gender
// Formula based on research averages: men typically have stride length of ~0.415 of height
// women typically have stride length of ~0.413 of height
func CalculateStrideLength(heightValue float64, heightUnit string, gender string) float64 {
	// Convert height to meters if in inches
	height := heightValue
	if heightUnit == "in" {
		height = height * 0.0254 // Convert inches to meters
	} else {
		height = height * 0.01 // Convert cm to meters
	}

	// Calculate stride length based on gender
	strideFactor := 0.415 // Default to male
	if gender == "female" {
		strideFactor = 0.413
	}

	return height * strideFactor
}

// CalculateTargetBPM calculates the target BPM for music based on running pace
// paceInSeconds: time in seconds for 1 km or 1 mile
// paceUnit: "km" or "mile"
func CalculateTargetBPM(paceInSeconds int, paceUnit string, strideLength float64) int {
	var distanceInMeters float64
	if paceUnit == "mile" {
		distanceInMeters = 1609.34 // 1 mile in meters
	} else {
		distanceInMeters = 1000.0 // 1 km in meters
	}

	// Calculate steps needed for the distance
	stepsNeeded := distanceInMeters / strideLength

	// Calculate steps per minute (BPM)
	// Multiply by 60 to convert from per second to per minute
	stepsPerMinute := (stepsNeeded / float64(paceInSeconds)) * 60.0

	// Most runners take 2 steps per beat for comfortable running
	// So we divide the steps per minute by 2 to get the target music BPM
	targetBPM := stepsPerMinute / 2.0

	return int(targetBPM)
} 