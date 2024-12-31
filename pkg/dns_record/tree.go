package record

//
// type RRNode struct {
// 	Name     string // e.g .com
// 	Records  ResourceRecord
// 	Children []*RRNode
// }
//
// // func RecordsFor(name []string, root *RRNode) []ResourceRecord {
// //
// // }
//
// func recordsForRecursive(name []string, level int, node *RRNode) []ResourceRecord {
//
// 	if level > len(name) {
// 		return []ResourceRecord{}
// 	}
//
// 	group := name[len(name)-level]
//
// 	for _, node := range node.Children {
// 		if node.Name == group {
//       return
// 		}
// 	}
// }
