import 'package:flutter/foundation.dart';
import '../../data/models/trader_model.dart';
import '../../data/repositories/copytrading_repository.dart';

class CopytradingProvider extends ChangeNotifier {
  final CopytradingRepository _repository;

  CopytradingProvider({CopytradingRepository? repository})
      : _repository = repository ?? CopytradingRepository();

  List<Trader> _traders = [];
  bool _isLoading = false;
  String? _error;

  List<Trader> get traders => _traders;
  bool get isLoading => _isLoading;
  String? get error => _error;

  Future<void> loadTraders() async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      _traders = await _repository.getTraders();
    } catch (e) {
      _error = e.toString();
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }
}
