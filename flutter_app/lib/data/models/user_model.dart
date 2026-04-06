import 'package:equatable/equatable.dart';
import 'package:json_annotation/json_annotation.dart';

part 'user_model.g.dart';

@JsonSerializable()
class UserModel extends Equatable {
  final int id;
  final String email;
  @JsonKey(name: 'api_key_set')
  final bool apiKeySet;

  const UserModel({
    required this.id,
    required this.email,
    required this.apiKeySet,
  });

  factory UserModel.fromJson(Map<String, dynamic> json) => _$UserModelFromJson(json);
  Map<String, dynamic> toJson() => _$UserModelToJson(this);

  @override
  List<Object?> get props => [id, email, apiKeySet];
}
