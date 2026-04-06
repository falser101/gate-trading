import '../../core/network/api_client.dart';
import '../../core/constants/api_constants.dart';
import '../models/balance_model.dart';
import '../models/account_detail_model.dart';

class AccountRepository {
  final ApiClient _client;

  AccountRepository(this._client);

  Future<TotalBalanceModel> getBalance() async {
    final response = await _client.get(ApiConstants.accountBalance);
    final data = response.data as Map<String, dynamic>;
    final balanceData = data['data'] as Map<String, dynamic>;
    return TotalBalanceModel.fromJson(balanceData);
  }

  Future<AccountDetailModel?> getAccountDetail() async {
    final response = await _client.get(ApiConstants.accountDetail);
    final data = response.data as Map<String, dynamic>;
    final detailData = data['data'] as Map<String, dynamic>;
    return AccountDetailModel.fromJson(detailData);
  }
}
