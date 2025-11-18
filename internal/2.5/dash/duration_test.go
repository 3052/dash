package dash

import (
   "testing"
)

func TestDuration_AsSeconds(t *testing.T) {
   testCases := []struct {
      name         string
      period       *Period
      expectedSecs float64
      expectErr    bool
   }{
      {"Simple Seconds", &Period{Duration: "PT10S"}, 10.0, false},
      {"Decimal Seconds", &Period{Duration: "PT3.5S"}, 3.5, false},
      {"Simple Minutes", &Period{Duration: "PT2M"}, 120.0, false},
      {"Simple Hours", &Period{Duration: "PT1H"}, 3600.0, false},
      {"Complex", &Period{Duration: "PT1H30M10.5S"}, 5410.5, false},
      {"Minutes and Seconds", &Period{Duration: "PT1M30S"}, 90.0, false},
      {"Hours and Seconds", &Period{Duration: "PT1H15S"}, 3615.0, false},
      {"Invalid Prefix", &Period{Duration: "P10S"}, 0, true},
      {"Unsupported Day", &Period{Duration: "P1DT12H"}, 0, true},
      {"Malformed", &Period{Duration: "PT10X"}, 0, true},
   }

   for _, tc := range testCases {
      t.Run(tc.name, func(t *testing.T) {
         secs, err := tc.period.AsSeconds()
         if tc.expectErr {
            if err == nil {
               t.Errorf("expected an error but got none")
            }
         } else {
            if err != nil {
               t.Errorf("did not expect an error but got: %v", err)
            }
            if secs != tc.expectedSecs {
               t.Errorf("expected %f seconds, but got %f", tc.expectedSecs, secs)
            }
         }
      })
   }
}
