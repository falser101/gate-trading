import 'package:flutter/foundation.dart';
import '../../core/network/api_client.dart';
import '../../data/repositories/strategy_repository.dart';
import '../../data/models/strategy_model.dart';

class StrategiesState {
  final bool isLoading;
  final List<StrategyModel> strategies;
  final String? error;

  StrategiesState({
    this.isLoading = false,
    this.strategies = const [],
    this.error,
  });

  StrategiesState copyWith({
    bool? isLoading,
    List<StrategyModel>? strategies,
    String? error,
  }) {
    return StrategiesState(
      isLoading: isLoading ?? this.isLoading,
      strategies: strategies ?? this.strategies,
      error: error,
    );
  }
}

class StrategyProvider extends ChangeNotifier {
  final StrategyRepository _repository;
  StrategiesState _state = StrategiesState();

  StrategyProvider({ApiClient? apiClient})
      : _repository = StrategyRepository(apiClient ?? ApiClient());

  StrategiesState get state => _state;
  List<StrategyModel> get strategies => _state.strategies;
  bool get isLoading => _state.isLoading;
  String? get error => _state.error;

  Future<void> loadStrategies() async {
    _state = _state.copyWith(isLoading: true, error: null);
    notifyListeners();
    try {
      final strategies = await _repository.getStrategies();
      _state = _state.copyWith(isLoading: false, strategies: strategies);
      notifyListeners();
    } catch (e) {
      _state = _state.copyWith(isLoading: false, error: e.toString());
      notifyListeners();
    }
  }

  Future<bool> createStrategy({
    required String type,
    required String symbol,
    required Map<String, dynamic> config,
  }) async {
    try {
      await _repository.createStrategy(type: type, symbol: symbol, config: config);
      await loadStrategies();
      return true;
    } catch (e) {
      return false;
    }
  }

  Future<void> startStrategy(int id) async {
    await _repository.startStrategy(id);
    await loadStrategies();
  }

  Future<void> stopStrategy(int id) async {
    await _repository.stopStrategy(id);
    await loadStrategies();
  }

  Future<void> deleteStrategy(int id) async {
    await _repository.deleteStrategy(id);
    await loadStrategies();
  }
}
