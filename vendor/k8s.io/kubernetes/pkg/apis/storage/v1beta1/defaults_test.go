/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1_test

import (
	"reflect"
	"testing"

	storagev1beta1 "k8s.io/api/storage/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	utilfeaturetesting "k8s.io/apiserver/pkg/util/feature/testing"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	_ "k8s.io/kubernetes/pkg/apis/storage/install"
	"k8s.io/kubernetes/pkg/features"
)

func roundTrip(t *testing.T, obj runtime.Object) runtime.Object {
	codec := legacyscheme.Codecs.LegacyCodec(storagev1beta1.SchemeGroupVersion)
	data, err := runtime.Encode(codec, obj)
	if err != nil {
		t.Errorf("%v\n %#v", err, obj)
		return nil
	}
	obj2, err := runtime.Decode(codec, data)
	if err != nil {
		t.Errorf("%v\nData: %s\nSource: %#v", err, string(data), obj)
		return nil
	}
	obj3 := reflect.New(reflect.TypeOf(obj).Elem()).Interface().(runtime.Object)
	err = legacyscheme.Scheme.Convert(obj2, obj3, nil)
	if err != nil {
		t.Errorf("%v\nSource: %#v", err, obj2)
		return nil
	}
	return obj3
}

func TestSetDefaultVolumeBindingMode(t *testing.T) {
	class := &storagev1beta1.StorageClass{}

	// When feature gate is enabled, field should be defaulted
	defaultMode := storagev1beta1.VolumeBindingImmediate
	output := roundTrip(t, runtime.Object(class)).(*storagev1beta1.StorageClass)
	outMode := output.VolumeBindingMode
	if outMode == nil {
		t.Errorf("Expected VolumeBindingMode to be defaulted to: %+v, got: nil", defaultMode)
	} else if *outMode != defaultMode {
		t.Errorf("Expected VolumeBindingMode to be defaulted to: %+v, got: %+v", defaultMode, outMode)
	}

	class = &storagev1beta1.StorageClass{}

	// When feature gate is disabled, field should not be defaulted
	defer utilfeaturetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, features.VolumeScheduling, false)()
	output = roundTrip(t, runtime.Object(class)).(*storagev1beta1.StorageClass)
	if output.VolumeBindingMode != nil {
		t.Errorf("Expected VolumeBindingMode to not be defaulted, got: %+v", output.VolumeBindingMode)
	}

}