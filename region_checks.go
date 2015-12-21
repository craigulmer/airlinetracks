
package airlinetracks




func (t Track) AnyPointsOutsideOfRegion2D(lleft Point, uright Point) (bool){
	for _,v := range t.Points {
		if (lleft.X > v.X) || (uright.X < v.X) ||
			(lleft.Y > v.Y) || (uright.Y < v.Y) {
			return true		
		}
	}
	return false
}

func (t Track) AnyPointsInsideRegion2D(lleft Point, uright Point) (bool){
	for _,v := range t.Points {
		if (lleft.X <= v.X) && (uright.X >= v.X) &&
			(lleft.Y <= v.Y) && (uright.Y >= v.Y) {
			return true		
		}
	}
	return false
}

func (t Track) AnyPointsInsideRegion3D(lleft Point, uright Point) (bool){
	for _,v := range t.Points {
		if (lleft.X <= v.X) && (uright.X >= v.X) &&
			(lleft.Y <= v.Y) && (uright.Y >= v.Y) &&
			(lleft.Z <= v.Z) && (uright.Z >= v.Z) {
			return true
		}
	}
   return false
}

func GetNamedRegionBoundingBox(name string) (Point, Point) {
	switch name {
	case "cuba"      : return Point{X:-84.957428,Y:19.828079}, Point{X:-74.131783,Y:23.283779}
	case "usa-main"  : return Point{X:-124.98,   Y:23.65},     Point{X:-66.88,    Y:49.21}
	default:
		panic("Unknow bounding box "+name)
	}
	return Point{},Point{} //Not called

}
