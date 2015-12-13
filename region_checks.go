
package airlinetracks




func (t Track) AnyPointsOutsideOfRegion2D(lleft point, uright point) (bool){
	for _,v := range t.points {
		if (lleft.x > v.x) || (uright.x < v.x) ||
			(lleft.y > v.y) || (uright.y < v.y) {
			return true		
		}
	}
	return false
}

func (t Track) AnyPointsInsideRegion2D(lleft point, uright point) (bool){
	for _,v := range t.points {
		if (lleft.x <= v.x) && (uright.x >= v.x) &&
			(lleft.y <= v.y) && (uright.y >= v.y) {
			return true		
		}
	}
	return false
}

func (t Track) AnyPointsInsideRegion3D(lleft point, uright point) (bool){
	for _,v := range t.points {
		if (lleft.x <= v.x) && (uright.x >= v.x) &&
			(lleft.y <= v.y) && (uright.y >= v.y) &&
			(lleft.z <= v.z) && (uright.z >= v.z) {
			return true
		}
	}
   return false
}

func GetNamedRegionBoundingBox(name string) (point, point) {
	switch name {
	case "cuba"      : return point{x:-84.957428,y:19.828079}, point{x:-74.131783,y:23.283779}
	case "usa-main"  : return point{x:-124.98,   y:23.65},     point{x:-66.88,    y:49.21}
	default:
		panic("Unknow bounding box "+name)
	}
	return point{},point{} //Not called

}
