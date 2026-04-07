import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/trader_model.dart';

class CopytradingRepository {
  final String baseUrl;

  CopytradingRepository({this.baseUrl = 'http://localhost:8080'});

  Future<List<Trader>> getTraders({
    int page = 1,
    int pageSize = 20,
  }) async {
    final url = Uri.parse(
      '$baseUrl/api/copytrading/traders?page=$page&page_size=$pageSize',
    );

    final response = await http.get(url);

    if (response.statusCode == 200) {
      final data = json.decode(response.body) as List;
      return data.map((e) => Trader.fromJson(e)).toList();
    } else {
      throw Exception('Failed to load traders: ${response.statusCode}');
    }
  }
}
