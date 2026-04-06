// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'futures_order_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

FuturesOrderModel _$FuturesOrderModelFromJson(Map<String, dynamic> json) =>
    FuturesOrderModel(
      id: (json['id'] as num?)?.toInt(),
      contract: json['contract'] as String?,
      size: json['size'] as String?,
      price: json['price'] as String?,
      status: json['status'] as String?,
      left: json['left'] as String?,
      fillPrice: json['fillPrice'] as String?,
      createTime: (json['createTime'] as num?)?.toDouble(),
    );

Map<String, dynamic> _$FuturesOrderModelToJson(FuturesOrderModel instance) =>
    <String, dynamic>{
      'id': instance.id,
      'contract': instance.contract,
      'size': instance.size,
      'price': instance.price,
      'status': instance.status,
      'left': instance.left,
      'fillPrice': instance.fillPrice,
      'createTime': instance.createTime,
    };
