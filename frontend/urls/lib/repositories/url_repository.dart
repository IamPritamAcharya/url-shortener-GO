import '../services/url_service.dart';
import '../models/url_model.dart';
import '../models/api_response.dart';

class UrlRepository {
  final UrlService _urlService;

  UrlRepository({required UrlService urlService}) : _urlService = urlService;

  Future<ApiResponse<UrlModel>> shortenUrl(String url, {String? customCode}) async {
    final request = ShortenRequest(url: url, customCode: customCode);
    return await _urlService.shortenUrl(request);
  }

  Future<ApiResponse<UrlModel>> getUrlStats(String code) async {
    return await _urlService.getUrlStats(code);
  }
}