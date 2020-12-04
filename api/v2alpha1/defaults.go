/*
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

// defaults.go - methods to get values for optional spec fields that return a default value when none set
//             - methods to initialize spec fields with default or derived values

package v2alpha1

import (
	"time"

	"github.com/iter8-tools/etc3/configuration"
)

const (
	// DefaultRollbackHandler is default name of rollback handler
	DefaultRollbackHandler string = "rollback"

	// DefaultMaxCandidateWeight is the default traffic percentage used in experiment, which is 100
	DefaultMaxCandidateWeight int32 = 100

	// DefaultMaxCandidateWeightIncrement is the default maxIncrement for traffic update, which is 10
	DefaultMaxCandidateWeightIncrement int32 = 10

	// DefaultBlueGreenAlgorithm is the default traffic weight recommendation algorithm for bluegreen experiments
	DefaultBlueGreenAlgorithm AlgorithmType = "FixedSplit"

	// DefaultAlgorithm is the default traffic weight recommendation algorithm for experiments other than bluegreen
	DefaultAlgorithm AlgorithmType = "Progressive"

	// DefaultIntervalSeconds is default interval duration as a string
	DefaultIntervalSeconds = 20

	// DefaultMaxIterations is the default number of iterations, which is 15
	DefaultMaxIterations int32 = 15
)

// DefaultBlueGreenSplit is the default split to be used for bluegreen experiment
var DefaultBlueGreenSplit = []int32{0, 100}

// GetNumberOfCandidates returns the number of candidates in VersionInfo
func (s *ExperimentSpec) GetNumberOfCandidates() int {
	if s.VersionInfo == nil {
		return 0
	}
	return len((*s.VersionInfo).Candidates)
}

// HasBaseline determines if a baseline has been identified in a s.VersionInfo
func (s *ExperimentSpec) HasBaseline() bool {
	return !(s.VersionInfo == nil)
}

// GetNumberOfBaseline returns the number of baselines in VersionInfo (1 if present, 0 otherwise)
func (s *ExperimentSpec) GetNumberOfBaseline() int {
	if s.HasBaseline() {
		return 1
	}
	return 0
}

//////////////////////////////////////////////////////////////////////
// spec.strategy.handlers
//////////////////////////////////////////////////////////////////////

func handlersForStrategy(cfg configuration.Iter8Config, strategy StrategyType) *configuration.Handlers {
	for _, t := range cfg.ExperimentTypes {
		if t.Name == string(strategy) {
			return &t.Handlers
		}
	}
	return nil
}

// GetStartHandler returns the name of the handler to be called when an experiment starts
func (s *ExperimentSpec) GetStartHandler(cfg configuration.Iter8Config) *string {
	if s.Strategy.Handlers == nil || s.Strategy.Handlers.Start == nil {
		handlers := handlersForStrategy(cfg, s.Strategy.Type)
		if nil == handlers {
			return nil
		}
		return &handlers.Start
	}
	return s.Strategy.Handlers.Start
}

// InitializeStartHandler iinitializes the start handler (if not already set) to the
// default rollback handler defined by the iter8 config.
func (s *ExperimentSpec) InitializeStartHandler(cfg configuration.Iter8Config) bool {
	if s.Strategy.Handlers == nil {
		s.Strategy.Handlers = &Handlers{}
	}
	if s.Strategy.Handlers.Start == nil {
		handler := s.GetStartHandler(cfg)
		if handler != nil {
			s.Strategy.Handlers.Start = handler
			return true
		}
	}
	return false
}

// GetFinishHandler returns the handler that should be called when an experiment ha completed.
func (s *ExperimentSpec) GetFinishHandler(cfg configuration.Iter8Config) *string {
	if s.Strategy.Handlers == nil || s.Strategy.Handlers.Finish == nil {
		handlers := handlersForStrategy(cfg, s.Strategy.Type)
		if nil == handlers {
			return nil
		}
		return &handlers.Finish
	}
	return s.Strategy.Handlers.Finish
}

// InitializeFinishHandler iinitializes the finish handler (if not already set) to the
// default rollback handler defined by the iter8 config.
func (s *ExperimentSpec) InitializeFinishHandler(cfg configuration.Iter8Config) bool {
	if s.Strategy.Handlers == nil {
		s.Strategy.Handlers = &Handlers{}
	}
	if s.Strategy.Handlers.Finish == nil {
		handler := s.GetFinishHandler(cfg)
		if handler != nil {
			s.Strategy.Handlers.Finish = handler
			return true
		}
	}
	return false
}

// GetRollbackHandler returns the handler to be called if a candidate fails its objective(s)
func (s *ExperimentSpec) GetRollbackHandler(cfg configuration.Iter8Config) *string {
	if s.Strategy.Handlers == nil || s.Strategy.Handlers.Rollback == nil {
		handlers := handlersForStrategy(cfg, s.Strategy.Type)
		if nil == handlers {
			return nil
		}
		return &handlers.Rollback
	}
	return s.Strategy.Handlers.Rollback
}

// InitializeRollbackHandler initializes the rollback handler (if not already set) to the
// default rollback handler defined by the iter8 config.
// Returns true if a change was made, false if not
func (s *ExperimentSpec) InitializeRollbackHandler(cfg configuration.Iter8Config) bool {
	if s.Strategy.Handlers == nil {
		s.Strategy.Handlers = &Handlers{}
	}
	if s.Strategy.Handlers.Rollback == nil {
		handler := s.GetRollbackHandler(cfg)
		if handler != nil {
			s.Strategy.Handlers.Rollback = handler
			return true
		}
	}
	return false
}

// GetFailureHandler returns the handler to be called if there is a failure during experiment execution
func (s *ExperimentSpec) GetFailureHandler(cfg configuration.Iter8Config) *string {
	if s.Strategy.Handlers == nil || s.Strategy.Handlers.Failure == nil {
		handlers := handlersForStrategy(cfg, s.Strategy.Type)
		if nil == handlers {
			return nil
		}
		return &handlers.Failure

	}
	return s.Strategy.Handlers.Failure
}

// InitializeFailureHandler initializes the finish handler (if not already set) to the default handler
// Returns true if a change was made, false if not
func (s *ExperimentSpec) InitializeFailureHandler(cfg configuration.Iter8Config) bool {
	if s.Strategy.Handlers == nil {
		s.Strategy.Handlers = &Handlers{}
	}
	if s.Strategy.Handlers.Failure == nil {
		handler := s.GetFailureHandler(cfg)
		if handler != nil {
			s.Strategy.Handlers.Failure = handler
			return true
		}
	}
	return false
}

// InitializeHandlers initialize handlers if not already set
// Returns true if a change was made, false if not
func (s *ExperimentSpec) InitializeHandlers(cfg configuration.Iter8Config) bool {
	changed := s.InitializeStartHandler(cfg)
	changed = s.InitializeFinishHandler(cfg) || changed
	changed = s.InitializeRollbackHandler(cfg) || changed
	changed = s.InitializeFailureHandler(cfg) || changed
	return changed
}

//////////////////////////////////////////////////////////////////////
// spec.strategy.weights
//////////////////////////////////////////////////////////////////////

// GetMaxCandidateWeight return spec.strategy.weights.maxCandidateWeight if set
// Otherwise it returns DefaultMaxCandidateWeight (100)
func (s *ExperimentSpec) GetMaxCandidateWeight() int32 {
	if s.Strategy.Weights == nil || s.Strategy.Weights.MaxCandidateWeight == nil {
		return DefaultMaxCandidateWeight
	}
	return *s.Strategy.Weights.MaxCandidateWeight
}

// InitializeMaxCandidateWeight initializes spec.strategy.weights.maxCandiateWeight if not already set
// Returns true if a change was made, false if not
func (s *ExperimentSpec) InitializeMaxCandidateWeight() bool {
	if s.Strategy.Weights == nil {
		s.Strategy.Weights = &Weights{}
	}
	if s.Strategy.Weights.MaxCandidateWeight == nil {
		weight := s.GetMaxCandidateWeight()
		s.Strategy.Weights.MaxCandidateWeight = &weight
		return true
	}
	return false
}

// GetMaxCandidateWeightIncrement return spec.strategy.weights.maxCandidateWeightIncrement if set
// Otherwise it returns DefaultMaxCandidateWeightIncrement (10)
func (s *ExperimentSpec) GetMaxCandidateWeightIncrement() int32 {
	if s.Strategy.Weights == nil || s.Strategy.Weights.MaxCandidateWeightIncrement == nil {
		return DefaultMaxCandidateWeightIncrement
	}
	return *s.Strategy.Weights.MaxCandidateWeightIncrement
}

// InitializeMaxCandidateWeightIncrement initializes spec.strategy.weights.maxCandidateWeightIncrement if not already set
// Returns true if a change was made, false if not
func (s *ExperimentSpec) InitializeMaxCandidateWeightIncrement() bool {
	if s.Strategy.Weights == nil {
		s.Strategy.Weights = &Weights{}
	}
	if s.Strategy.Weights.MaxCandidateWeightIncrement == nil {
		increment := s.GetMaxCandidateWeightIncrement()
		s.Strategy.Weights.MaxCandidateWeightIncrement = &increment
		return true
	}
	return false
}

// GetAlgorithm returns spec.strategy.weights.algorithm if set
// Otherwise it returns the default based on spec.strategy.type
func (s *ExperimentSpec) GetAlgorithm() AlgorithmType {
	if s.Strategy.Weights == nil || s.Strategy.Weights.Algorithm == nil {
		switch s.Strategy.Type {
		case StrategyTypeBlueGreen:
			return DefaultBlueGreenAlgorithm
		}
		return DefaultAlgorithm
	}
	return *s.Strategy.Weights.Algorithm
}

// InitializeAlgorithm initializes spec.strategy.weights.algorithm if not already set
// Returns true if a change was made, false if not
func (s *ExperimentSpec) InitializeAlgorithm() bool {
	if s.Strategy.Weights == nil {
		s.Strategy.Weights = &Weights{}
	}
	if s.Strategy.Weights.Algorithm == nil {
		algorithm := s.GetAlgorithm()
		s.Strategy.Weights.Algorithm = &algorithm
		return true
	}
	return false
}

// UniformSplit returns the default (uniform) split for non-bluegreen experiments
func UniformSplit(numberOfCandidates int, maxCandidateWeight int32) []int32 {
	numCandidates := int32(numberOfCandidates)
	split := make([]int32, numberOfCandidates+1)
	if len(split) == 0 {
		return split
	}
	weight := maxCandidateWeight / numCandidates
	// candidate will get any "extra" caused by rounding
	split[0] = 100 - numCandidates*weight
	for i := 1; i <= int(numberOfCandidates); i++ {
		split[i] = weight
	}
	return split
}

// GetWeightDistribution returns spec.strategy.weights.split if set
// Otherwise it returns the default based on spec.strategy.type
func (s *ExperimentSpec) GetWeightDistribution() []int32 {
	if s.Strategy.Weights == nil || s.Strategy.Weights.WeightDistribution == nil {
		switch s.Strategy.Type {
		case StrategyTypeBlueGreen:
			// we expect Algorithm to be AlgorithmTypeFixedSplit
			return DefaultBlueGreenSplit
		}
		switch s.GetAlgorithm() {
		case AlgorithmTypeFixedSplit:
			return UniformSplit(s.GetNumberOfCandidates(), s.GetMaxCandidateWeight())
		}
		return make([]int32, 0)
	}
	return s.Strategy.Weights.WeightDistribution
}

// InitializeWeightDistribution initializes spec.strategy.weights.split if not already set (and algorithm is "fixed_split")
// Returns true if a change was made, false if not
func (s *ExperimentSpec) InitializeWeightDistribution() bool {
	if s.Strategy.Weights == nil {
		s.Strategy.Weights = &Weights{}
	}
	if s.Strategy.Weights.WeightDistribution == nil {
		s.Strategy.Weights.WeightDistribution = s.GetWeightDistribution()
		return true
	}
	return false
}

// InitializeWeights initializes weights if not already set
// Returns true if a change was made, false if not
func (s *ExperimentSpec) InitializeWeights() bool {
	change := s.InitializeMaxCandidateWeight()
	change = s.InitializeMaxCandidateWeightIncrement() || change
	change = s.InitializeAlgorithm() || change
	// Must wait until versionInfo has been defined by start handler before
	// initializing weight distribution because need to know the candidates
	// change = s.InitializeWeightDistribution() || change
	return change
}

//////////////////////////////////////////////////////////////////////
// spec.duration
//////////////////////////////////////////////////////////////////////

// GetIntervalSeconds returns specified(or default) interval for each duration
func (s *ExperimentSpec) GetIntervalSeconds() int32 {
	if s.Duration == nil || s.Duration.IntervalSeconds == nil {
		return DefaultIntervalSeconds
	}
	return *s.Duration.IntervalSeconds
}

// GetIntervalAsDuration returns spec.duration.intervalSeconds as a time.Duration (in ns)
func (s *ExperimentSpec) GetIntervalAsDuration() time.Duration {
	return time.Second * time.Duration(s.GetIntervalSeconds())
}

// InitializeInterval sets duration.interval if not already set using the default value
// Returns true if a change was made, false if not
func (s *ExperimentSpec) InitializeInterval() bool {
	if s.Duration == nil {
		s.Duration = &Duration{}
	}
	if s.Duration.IntervalSeconds == nil {
		interval := int32(DefaultIntervalSeconds)
		s.Duration.IntervalSeconds = &interval
		return true
	}
	return false
}

// GetMaxIterations returns specified(or default) max of iterations
func (s *ExperimentSpec) GetMaxIterations() int32 {
	if s.Duration == nil || s.Duration.MaxIterations == nil {
		return DefaultMaxIterations
	}
	return *s.Duration.MaxIterations
}

// InitializeMaxIterations sets duration.maxIterations if not already set using the default value
// Returns true if a change was made, false if not
func (s *ExperimentSpec) InitializeMaxIterations() bool {
	if s.Duration == nil {
		s.Duration = &Duration{}
	}
	if s.Duration.MaxIterations == nil {
		iterations := s.GetMaxIterations()
		s.Duration.MaxIterations = &iterations
		return true
	}
	return false
}

// InitializeDuration initializes spec.durations if not already set
// Returns true if a change was made, false if not
func (s *ExperimentSpec) InitializeDuration() bool {
	change := s.InitializeInterval()
	change = s.InitializeMaxIterations() || change
	return change
}

//////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////
// criteria
//////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////

// GetRequestCount returns the requst count metric
// If there are no criteria specified, this is nil
func (s *ExperimentSpec) GetRequestCount(cfg configuration.Iter8Config) *string {
	if s.Criteria == nil {
		return nil
	}
	if s.Criteria.RequestCount == nil {
		rc := cfg.RequestCount
		return &rc
	}
	return s.Criteria.RequestCount
}

// InitializeRequestCount sets the request count metric to the default value if not already set
// Returns true if a change was made, false if not
func (s *ExperimentSpec) InitializeRequestCount(cfg configuration.Iter8Config) bool {
	if s.Criteria == nil {
		return false
	}
	if s.Criteria.RequestCount == nil {
		s.Criteria.RequestCount = s.GetRequestCount(cfg)
		return true
	}
	return false
}

//////////////////////////////////////////////////////////////////////
// objective
//////////////////////////////////////////////////////////////////////

// GetRollbackOnFailure identifies if the experiment should be rolledback on failure of an objective
func (o *Objective) GetRollbackOnFailure(strategy StrategyType) bool {
	if o.RollbackOnFailure == nil {
		if strategy == StrategyTypeBlueGreen {
			return true
		}
		return false
	}
	return *o.RollbackOnFailure
}

// InitializeObjectives initializes the rollbackOnFailure field of all objectives if
// the strategy type is "bluegreen"
// Returns true if a change was made, false if not
func (s *ExperimentSpec) InitializeObjectives() bool {
	if s.Criteria == nil {
		return false
	}

	change := false
	for _, o := range s.Criteria.Objectives {
		if s.Strategy.Type == StrategyTypeBlueGreen && o.RollbackOnFailure == nil {
			rollback := true
			o.RollbackOnFailure = &rollback
			change = true
		}
	}
	return change
}

// InitializeCriteria initializes any criteria details not already set
// Returns true if a change was made, false if not
func (s *ExperimentSpec) InitializeCriteria(cfg configuration.Iter8Config) bool {
	if s.Criteria == nil {
		return false
	}
	change := s.InitializeRequestCount(cfg)
	change = s.InitializeObjectives() || change
	return change
}