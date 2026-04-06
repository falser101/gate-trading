import '../../core/network/api_client.dart';
import '../../core/constants/api_constants.dart';
import '../models/user_model.dart';

class AuthRepository {
  final ApiClient _client;

  AuthRepository(this._client);

  ApiClient get apiClient => _client;

  Future<Map<String, dynamic>> register(String email, String password) async {
    final response = await _client.post(
      ApiConstants.register,
      data: {'email': email, 'password': password},
    );
    return response.data as Map<String, dynamic>;
  }

  Future<Map<String, dynamic>> login(String email, String password) async {
    final response = await _client.post(
      ApiConstants.login,
      data: {'email': email, 'password': password},
    );
    final data = response.data as Map<String, dynamic>;
    if (data['token'] != null) {
      await _client.setToken(data['token']);
    }
    return data;
  }

  Future<void> logout() async {
    await _client.clearToken();
  }

  Future<UserModel?> getUser() async {
    try {
      final response = await _client.get(ApiConstants.user);
      return UserModel.fromJson(response.data as Map<String, dynamic>);
    } catch (e) {
      return null;
    }
  }

  Future<void> bindApiKey(String apiKey, String apiSecret) async {
    await _client.post(
      ApiConstants.apiKey,
      data: {'api_key': apiKey, 'api_secret': apiSecret},
    );
  }
}
