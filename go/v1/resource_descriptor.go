/*
Wrapper APIs for in-toto attestation ResourceDescriptor protos.
*/

package v1

func (d *ResourceDescriptor) Validate() bool {
	// at least one of name, URI or digest are required
	if d.GetName() == "" && d.GetUri() == "" && len(d.GetDigest()) == 0 {
		return false
	}

	return true
}
