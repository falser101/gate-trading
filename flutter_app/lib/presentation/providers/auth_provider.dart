import 'package:flutter/foundation.dart';
import '../../core/network/api_client.dart';
import '../../data/repositories/auth_repository.dart';
import '../../data/models/user_model.dart';
import '../../core/router.dart';

class AuthState {
  final bool isLoading;
  final bool isAuthenticated;
  final UserModel? user;
  final String? error;

  AuthState({
    this.isLoading = false,
    this.isAuthenticated = false,
    this.user,
    this.error,
  });

  AuthState copyWith({
    bool? isLoading,
    bool? isAuthenticated,
    UserModel? user,
    String? error,
  }) {
    return AuthState(
      isLoading: isLoading ?? this.isLoading,
      isAuthenticated: isAuthenticated ?? this.isAuthenticated,
      user: user ?? this.user,
      error: error,
    );
  }
}

class AuthProvider extends ChangeNotifier {
  final AuthRepository _repository;
  AuthState _state = AuthState();

  AuthProvider({ApiClient? apiClient})
      : _repository = AuthRepository(apiClient ?? ApiClient()) {
    _checkAuth();
  }

  AuthState get state => _state;
  UserModel? get user => _state.user;
  bool get isAuthenticated => _state.isAuthenticated;
  bool get isLoading => _state.isLoading;
  String? get error => _state.error;

  Future<void> _checkAuth() async {
    final token = await _repository.apiClient.getToken();
    if (token != null) {
      final user = await _repository.getUser();
      if (user != null) {
        _state = _state.copyWith(isAuthenticated: true, user: user);
        updateAuthState(true);
        notifyListeners();
      }
    }
  }

  Future<bool> login(String email, String password) async {
    _state = _state.copyWith(isLoading: true, error: null);
    notifyListeners();
    try {
      await _repository.login(email, password);
      final user = await _repository.getUser();
      _state = _state.copyWith(
        isLoading: false,
        isAuthenticated: true,
        user: user,
      );
      updateAuthState(true);
      notifyListeners();
      return true;
    } catch (e) {
      _state = _state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      notifyListeners();
      return false;
    }
  }

  Future<bool> register(String email, String password) async {
    _state = _state.copyWith(isLoading: true, error: null);
    notifyListeners();
    try {
      await _repository.register(email, password);
      final user = await _repository.getUser();
      _state = _state.copyWith(
        isLoading: false,
        isAuthenticated: true,
        user: user,
      );
      updateAuthState(true);
      notifyListeners();
      return true;
    } catch (e) {
      _state = _state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      notifyListeners();
      return false;
    }
  }

  Future<void> logout() async {
    await _repository.logout();
    _state = AuthState();
    updateAuthState(false);
    notifyListeners();
  }

  Future<void> bindApiKey(String apiKey, String apiSecret) async {
    await _repository.bindApiKey(apiKey, apiSecret);
    final user = await _repository.getUser();
    if (user != null) {
      _state = _state.copyWith(user: user);
      notifyListeners();
    }
  }
}
