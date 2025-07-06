import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/url_model.dart';
import '../models/api_response.dart';

class UrlService {
  static const String baseUrl = 'http://localhost:8080';

  Future<ApiResponse<UrlModel>> shortenUrl(ShortenRequest request) async {
    try {
      final endpoint = request.customCode != null ? '/shorten/custom' : '/shorten';
      
      final response = await http.post(
        Uri.parse('$baseUrl$endpoint'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode(request.toJson()),
      );

      final data = json.decode(response.body);

      if (response.statusCode == 201) {
        return ApiResponse.success(UrlModel.fromJson(data));
      } else {
        return ApiResponse.error(
          data['message'] ?? 'Failed to shorten URL',
          statusCode: response.statusCode,
        );
      }
    } catch (e) {
      return ApiResponse.error('Network error: ${e.toString()}');
    }
  }

  Future<ApiResponse<UrlModel>> getUrlStats(String code) async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/stats/$code'),
        headers: {'Content-Type': 'application/json'},
      );

      final data = json.decode(response.body);

      if (response.statusCode == 200) {
        return ApiResponse.success(UrlModel.fromJson(data));
      } else {
        return ApiResponse.error(
          data['message'] ?? 'Failed to get statistics',
          statusCode: response.statusCode,
        );
      }
    } catch (e) {
      return ApiResponse.error('Network error: ${e.toString()}');
    }
  }
}